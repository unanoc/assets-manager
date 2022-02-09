package events

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	gh "github.com/google/go-github/v38/github"
	log "github.com/sirupsen/logrus"

	assetsmanager "github.com/trustwallet/assets-go-libs/client/assets-manager"
	"github.com/trustwallet/assets-go-libs/file"
	"github.com/trustwallet/assets-go-libs/http"
	"github.com/trustwallet/assets-go-libs/image"
	"github.com/trustwallet/assets-go-libs/path"
	"github.com/trustwallet/assets-go-libs/validation"
	"github.com/trustwallet/assets-go-libs/validation/list"
	"github.com/trustwallet/assets-manager/internal/config"
	"github.com/trustwallet/assets-manager/internal/services/worker/blockchain"
	"github.com/trustwallet/assets-manager/internal/services/worker/github"
	"github.com/trustwallet/assets-manager/internal/services/worker/metrics"
	"github.com/trustwallet/go-primitives/coin"
	"github.com/trustwallet/go-primitives/types"
)

type Handler struct {
	metrics       *metrics.Prometheus
	github        *github.Client
	blockchain    *blockchain.Client
	assetsManager *assetsmanager.Client
}

func NewHandler(metricsClient *metrics.Prometheus, githubClient *github.Client,
	blockchainClient *blockchain.Client, assetsManager *assetsmanager.Client) *Handler {
	return &Handler{
		metrics:       metricsClient,
		github:        githubClient,
		blockchain:    blockchainClient,
		assetsManager: assetsManager,
	}
}

func (e Handler) HandlePullRequestOpened(ctx context.Context, event *gh.PullRequestEvent) error {
	e.metrics.IncCounterPullRequestsCreated()

	owner := event.GetRepo().GetOwner().GetLogin()
	repo := event.GetRepo().GetName()
	prNum := event.GetPullRequest().GetNumber()
	prCreator := event.GetPullRequest().GetUser().GetLogin()

	log.WithFields(log.Fields{
		"pr_num":  prNum,
		"creator": prCreator,
	}).Debug("Pull request opened")

	if err := e.HandlePullRequestChangesPushed(ctx, event); err != nil {
		return err
	}

	if e.isCollaborator(prCreator) {
		return nil
	}

	if err := e.github.SetLabelOnPullRequest(ctx, owner, repo, prNum, &gh.Label{
		Name: gh.String(config.Default.Label.Requested),
	}); err != nil {
		return err
	}

	pp := getPaymentParams(event.GetPullRequest())
	commentText := substituteDynamicContent(config.Default.Message.Initial, &contentParams{PP: pp})

	return e.github.CreateCommentOnPullRequest(ctx, owner, repo, commentText, prNum)
}

func (e Handler) HandleIssueCommentCreated(ctx context.Context, event *gh.IssueCommentEvent) error {
	if !event.GetIssue().IsPullRequest() {
		return nil
	}

	owner := event.GetRepo().GetOwner().GetLogin()
	repo := event.GetRepo().GetName()
	commentBody := event.GetComment().GetBody()
	prCreator := event.GetIssue().GetUser().GetLogin()
	commentCreator := event.GetComment().GetUser().GetLogin()
	prNum := event.GetIssue().GetNumber()
	commentID := event.GetComment().GetID()

	log.WithFields(log.Fields{
		"pr_num":  prNum,
		"creator": commentCreator,
	}).Debug("Issued comment created")

	if e.isCollaborator(prCreator) {
		return nil
	}

	pr, err := e.github.GetPullRequest(ctx, owner, repo, prNum)
	if err != nil {
		return err
	}

	debugCheck := strings.Contains(commentBody, "/check")

	err = e.deleteCommentIfNeeded(ctx, owner, repo, prCreator, commentCreator, commentID)
	if err != nil {
		return err
	}

	return e.checkPullStatus(ctx, owner, repo, pr, debugCheck)
}

func (e Handler) HandlePullRequestReviewCommentCreated(ctx context.Context,
	event *gh.PullRequestReviewCommentEvent,
) error {
	owner := event.GetRepo().GetOwner().GetLogin()
	repo := event.GetRepo().GetName()
	prNum := event.GetPullRequest().GetNumber()
	commentCreator := event.GetComment().GetUser().GetLogin()

	log.WithFields(log.Fields{
		"pr_num":  prNum,
		"creator": commentCreator,
	}).Debug("Review comment created")

	pr, err := e.github.GetPullRequest(ctx, owner, repo, prNum)
	if err != nil {
		return err
	}

	return e.checkPullStatus(ctx, owner, repo, pr, false)
}

func (e Handler) deleteCommentIfNeeded(ctx context.Context, owner, repo, prCreator,
	user string, commentID int64,
) error {
	if !config.Default.UserAccess.DeleteCommentsFromExternal {
		return nil
	}

	if valid := e.isUserCollaboratorOrCreator(prCreator, user); valid {
		return nil
	}

	return e.github.DeleteCommentInIssue(ctx, owner, repo, commentID)
}

func (e Handler) isUserCollaboratorOrCreator(creator, user string) bool {
	isBot := strings.HasPrefix(user, config.Default.ServiceName)
	isCreator := user == creator

	return isBot || isCreator || e.isCollaborator(user)
}

func (e Handler) isCollaborator(user string) bool {
	isCollaborator := false

	if config.Default.UserAccess.Collaborators != "" {
		collaborators := strings.Split(config.Default.UserAccess.Collaborators, ",")
		for _, collaborator := range collaborators {
			if user == collaborator {
				isCollaborator = true
			}
		}
	}

	return isCollaborator
}

func (e Handler) checkPullStatus(ctx context.Context, owner, repo string, pr *gh.PullRequest, debug bool) error {
	if e.isCollaborator(pr.GetUser().GetLogin()) {
		return nil
	}

	if pr.GetState() != "open" {
		return nil
	}

	if !e.isPaymentExpected(ctx, owner, repo, pr) {
		if debug {
			text := substituteDynamicContent(config.Default.Message.Reviewed, nil)

			return e.github.CreateCommentOnPullRequest(ctx, owner, repo, text, pr.GetNumber())
		}

		return nil
	}

	// Check for already paid -> approve pr.
	paymentStatus, err := e.checkPaymentForPullRequest(pr)
	if err != nil {
		return err
	}
	if paymentStatus.Paid {
		return e.approvePullRequest(ctx, owner, repo, pr, paymentStatus)
	}

	now := time.Now()
	prAgeHours := now.Sub(pr.GetCreatedAt())
	prUpdateAgeHours := now.Sub(pr.GetUpdatedAt())
	halfHour := time.Duration(30) * time.Minute

	// Check for too old -> close pr.
	if prAgeHours >= config.Default.Timeout.MaxAgeClose && prUpdateAgeHours > halfHour {
		return e.closePullRequest(ctx, owner, repo, pr)
	}

	if debug {
		text := substituteDynamicContent(config.Default.Message.NotReceived, nil)

		return e.github.CreateCommentOnPullRequest(ctx, owner, repo, text, pr.GetNumber())
	}

	// Check if it's time for reminder.
	if prUpdateAgeHours >= config.Default.Timeout.MaxIdleRemind {
		return e.remindToPay(ctx, owner, repo, pr)
	}

	return nil
}

func (e Handler) isPaymentExpected(ctx context.Context, owner, repo string, pr *gh.PullRequest) bool {
	return !(e.hasReviewAlready(ctx, owner, repo, pr) || e.hasLabelAlready(ctx, owner, repo, pr))
}

func (e Handler) approvePullRequest(ctx context.Context, owner, repo string,
	pr *gh.PullRequest, ps *blockchain.PaymentStatus,
) error {
	text := substituteDynamicContent(config.Default.Message.Received, &contentParams{
		PaidAmount:       ps.Amount,
		PaidSymbol:       strings.Split(ps.Token, "-")[0],
		PaidExplorerLink: ps.Transactions[0].ExplorerLink,
		Moderators:       config.Default.UserAccess.Moderators,
	})

	if _, err := e.github.CreateReview(ctx, owner, repo, text, "APPROVE", pr.GetNumber()); err != nil {
		return err
	}

	if err := e.github.SetLabelOnPullRequest(ctx, owner, repo, pr.GetNumber(), &gh.Label{
		Name: gh.String(config.Default.Label.Paid),
	}); err != nil {
		return err
	}

	assignedUsers := strings.Split(config.Default.UserAccess.Moderators, ",")
	if _, err := e.github.AddAssignees(ctx, owner, repo, pr.GetNumber(), assignedUsers); err != nil {
		return err
	}

	e.metrics.IncCounterPaymentsDetected()

	explorerLink, err := e.blockchain.BurnToken(ps.Token, int64(ps.Amount*blockchain.AmountPrecision))
	if err != nil || explorerLink == "" {
		return err
	}

	text = substituteDynamicContent(config.Default.Message.Burned, &contentParams{
		PaidAmount:       ps.Amount,
		PaidSymbol:       strings.Split(ps.Token, "-")[0],
		BurnExplorerLink: explorerLink,
	})

	return e.github.CreateCommentOnPullRequest(ctx, owner, repo, text, pr.GetNumber())
}

func (e Handler) closePullRequest(ctx context.Context, owner, repo string, pr *gh.PullRequest) error {
	text := substituteDynamicContent(config.Default.Message.ClosingOldPR, nil)
	if err := e.github.CreateCommentOnPullRequest(ctx, owner, repo, text, pr.GetNumber()); err != nil {
		return err
	}

	return e.github.ClosePullRequest(ctx, owner, repo, pr.GetNumber())
}

func (e Handler) remindToPay(ctx context.Context, owner, repo string, pr *gh.PullRequest) error {
	pp := getPaymentParams(pr)
	text := substituteDynamicContent(config.Default.Message.Reminder, &contentParams{PP: pp})

	return e.github.CreateCommentOnPullRequest(ctx, owner, repo, text, pr.GetNumber())
}

func (e Handler) hasReviewAlready(ctx context.Context, owner, repo string, pr *gh.PullRequest) bool {
	list, err := e.github.GetPullRequestReviewList(ctx, owner, repo, pr.GetNumber())
	if err != nil {
		return false
	}

	for _, review := range list {
		if strings.HasPrefix(review.GetUser().GetLogin(), config.Default.ServiceName) {
			if review.GetState() == "APPROVED" {
				return true
			}
		}
	}

	return false
}

func (e Handler) hasLabelAlready(ctx context.Context, owner, repo string, pr *gh.PullRequest) bool {
	labels, err := e.github.GetIssueListLabels(ctx, owner, repo, pr.GetNumber())
	if err != nil {
		return false
	}

	for _, label := range labels {
		if label.GetName() == config.Default.Label.Paid {
			return true
		}
	}

	return false
}

func (e Handler) checkPaymentForPullRequest(pr *gh.PullRequest) (*blockchain.PaymentStatus, error) {
	params := getPaymentParams(pr)

	for _, p := range params.Payments {
		txs, err := e.blockchain.GetTransactionsForAddress(params.Address)
		if err != nil {
			return nil, err
		}

		ps := blockchain.GetPaymentStatus(
			txs, params.Address, p.Memo, p.Token, p.CreatedTime, p.EndTime, p.MinAmount)

		if ps.Paid {
			return ps, nil
		}
	}

	return &blockchain.PaymentStatus{}, nil
}

func (e Handler) CheckOpenPullRequests(ctx context.Context, owner, repo string, pr *gh.PullRequest) error {
	prs, err := e.github.GetPullRequestsList(ctx, owner, repo, "open", 100)
	if err != nil {
		return fmt.Errorf("failed to get open pull requests: %w", err)
	}

	var triggeredPrIncluded bool
	if pr != nil {
		for _, p := range prs {
			if p.GetNumber() == pr.GetNumber() {
				triggeredPrIncluded = true
			}
		}
	}

	if triggeredPrIncluded {
		prs = append(prs, pr)
	}

	e.metrics.SetPullRequestsOpen(len(prs))

	prCountToPay := 0
	for _, p := range prs {
		err := e.checkPullStatus(ctx, owner, repo, p, false)
		if err != nil {
			return err
		}

		if e.isPaymentExpected(ctx, owner, repo, pr) {
			prCountToPay++
		}
	}

	e.metrics.SetPullRequestsToPay(prCountToPay)

	return nil
}

func (e Handler) HandlePullRequestChangesPushed(ctx context.Context, event *gh.PullRequestEvent) error {
	owner := event.GetRepo().GetOwner().GetLogin()
	repo := event.GetRepo().GetName()
	pr := event.GetPullRequest()
	branch := event.GetPullRequest().GetHead().GetRef()
	headOwner := event.GetPullRequest().GetHead().GetRepo().GetOwner().GetLogin()
	headRepo := event.GetPullRequest().GetHead().GetRepo().GetName()

	log.WithFields(log.Fields{
		"pr_num":  pr.GetNumber(),
		"creator": headOwner,
	}).Debug("Pull request changes are pushed")

	files, err := e.github.GetPullRequestFileList(ctx, owner, repo, pr.GetNumber(), 100)
	if err != nil {
		return err
	}

	filesCheckSummary := e.getFilesCheckSummary(files, owner)
	tokenCheckSummary := e.getTokensCheckSummary(files, headOwner, headRepo, branch)
	validatorsCheckSummary := e.getValidatorsCheckSummary(files, headOwner, headRepo, branch)

	summary := fmt.Sprintf("%s%s\n%s", filesCheckSummary, tokenCheckSummary, validatorsCheckSummary)
	if tokenCheckSummary == "" && validatorsCheckSummary == "" {
		summary = "No token files found. If you try to add/modify a token, " +
			"check the name and location of your files! Logo file must be named exactly 'logo.png'. " +
			"If you are not adding a token, ignore this message."
	}

	err = e.github.CreateCommentOnPullRequest(ctx, owner, repo, summary, pr.GetNumber())
	if err != nil {
		return err
	}

	err = e.checkPullStatus(ctx, owner, repo, pr, false)
	if err != nil {
		return err
	}

	return nil
}

func (e Handler) getFilesCheckSummary(files []*gh.CommitFile, repoOwner string) string {
	text := "### PR Summary\n"

	checkSummary := e.checkPullRequestFiles(files, config.Default.Limitation.PrFilesNumMax, repoOwner)
	if checkSummary != "" {
		return fmt.Sprintf("%s%s", text, checkSummary)
	}

	return fmt.Sprintf("%sFiles OK: %d\n", text, len(files))
}

// nolint: gosec
func (e Handler) getTokensCheckSummary(files []*gh.CommitFile, repoOwner, repoName, branch string) string {
	// Check tokens.
	tokenIDs := make(map[string]string)

	for _, file := range files {
		id, tokenType := path.GetTokenFromAssetLogoPath(file.GetFilename())

		if id != "" && tokenType != "" {
			tokenIDs[id] = tokenType
		}
	}

	if len(tokenIDs) == 0 {
		return ""
	}

	var text string

	tokenHeaderTxt := "Token in PR: %s %s"
	if len(tokenIDs) > 1 {
		text += fmt.Sprintf("Tokens in PR: (%d)\n", len(tokenIDs))
		tokenHeaderTxt = "\n- %s %s"
	}

	for id, tokenType := range tokenIDs {
		text += fmt.Sprintf(tokenHeaderTxt, tokenType, id)
	}

	for id, tokenType := range tokenIDs {
		if len(tokenIDs) > 1 {
			text += fmt.Sprintf("\n-----\n**Token %s - %s**:", tokenType, id)
		}

		msg := e.checkToken(id, tokenType, repoOwner, repoName, branch)
		text += fmt.Sprintf("\n%s\n", msg)
	}

	return text
}

func (e Handler) getValidatorsCheckSummary(files []*gh.CommitFile, repoOwner, repoName, branch string) string {
	validatorLists := make([]*file.Path, 0)
	validatorAssetLogos := make([]*file.Path, 0)

	for _, f := range files {
		assetPath := file.NewPath(f.GetFilename())

		if assetPath.Type() == file.TypeValidatorsListFile {
			validatorLists = append(validatorLists, assetPath)
		}

		if assetPath.Type() == file.TypeValidatorsLogoFile {
			validatorAssetLogos = append(validatorAssetLogos, assetPath)
		}
	}

	text := "**Validators check**\n"
	if (len(validatorAssetLogos) == 0 && len(validatorLists) > 0) ||
		(len(validatorAssetLogos) > 0 && len(validatorLists) == 0) {
		return fmt.Sprintf("%s❌ For adding asset validators, you need to add validator logo and update list.json", text)
	}

	if len(validatorAssetLogos) == 0 && len(validatorLists) == 0 {
		return ""
	}

	errorsMsg := e.checkValidators(validatorLists, validatorAssetLogos, repoOwner, repoName, branch)
	if errorsMsg != "" {
		return fmt.Sprintf("%s%s", text, errorsMsg)
	}

	return text + "✅ Check OK"
}

func (e Handler) checkValidators(validatorLists, validatorAssetLogos []*file.Path,
	repoOwner, repoName, branch string) string {
	var errorsMsg string

	for _, vlist := range validatorLists {
		listURL := path.GetValidatorListGithubURL(repoOwner, repoName, branch, vlist.Chain().Handle)
		bytes, err := http.GetHTTPResponseBytes(listURL)
		if err != nil {
			return fmt.Sprintf("Failed to get file content of [list.json](%s)", listURL)
		}

		var validatorList []list.Model
		err = json.Unmarshal(bytes, &validatorList)
		if err != nil {
			return fmt.Sprintf("Failed to parse content of [list.json](%s)", listURL)
		}

		validatorMap := make(map[string]struct{})
		for _, v := range validatorList {
			validatorMap[*v.ID] = struct{}{}
		}

		for _, vlogo := range validatorAssetLogos {
			if vlogo.Chain() != vlist.Chain() {
				continue
			}

			if _, exists := validatorMap[vlogo.Asset()]; !exists {
				errorsMsg += fmt.Sprintf("❌ Asset '%s' (%s) not found in [list.json](%s)\n",
					vlogo.Asset(), vlogo.Chain().Handle, listURL)
			}

			logoURL := path.GetValidatorAssetLogoGithubURL(repoOwner, repoName, branch, vlogo.Chain().Handle, vlogo.Asset())
			errorsMsg += fmt.Sprintf("%s\n\n%s", e.checkLogo(logoURL), getLogoHTML(logoURL))

			errorsMsg += "\n-----\n"
		}
	}

	return errorsMsg
}

func (e Handler) checkPullRequestFiles(files []*gh.CommitFile, limit int, repoOwner string) string {
	if len(files) == 0 {
		return "No changed files found."
	}

	if len(files) > limit && !e.isCollaborator(repoOwner) {
		return fmt.Sprintf("Too many changed files: %d (max %d).", len(files), limit)
	}

	var msg string

	for _, file := range files {
		if err := validation.ValidateFileInPR(file.GetFilename()); err != nil {
			msg += fmt.Sprintf("%s: Please revert it.\n", err.Error())
		}

		if file.GetStatus() == "removed" {
			msg += fmt.Sprintf("File `%s` is being deleted. Files should not be deleted in a PR. "+
				"(Deprecated tokens should be deactivated only.)\n", file.GetFilename())
		}
	}

	return msg
}

func (e Handler) checkToken(tokenID, tokenType, repoOwner, repoName, branch string) string {
	chain, err := types.GetChainFromAssetType(tokenType)
	if err != nil {
		return "failed to get chain from asset type"
	}

	logoURL := path.GetAssetLogoGithubURL(repoOwner, repoName, branch, chain.Handle, tokenID)
	infoURL := path.GetAssetInfoGithubURL(repoOwner, repoName, branch, chain.Handle, tokenID)

	tokenInfo := &assetsmanager.AssetValidationReq{}
	err = http.GetHTTPResponse(infoURL, tokenInfo)
	if err != nil {
		return fmt.Sprintf("Failed to get info.json content: %s (%s)", err.Error(), infoURL)
	}

	var text string
	if tokenInfo.Symbol != nil {
		text += fmt.Sprintf("Symbol: **%s** ", *tokenInfo.Symbol)
	} else {
		text += "Symbol: **?** "
	}

	if tokenInfo.Decimals != nil {
		text += fmt.Sprintf("decimals: %d\n", *tokenInfo.Decimals)
	}

	text += fmt.Sprintf("Links: [Logo](%s) [Info](%s) ", logoURL, infoURL)

	if tokenInfo.Website != nil {
		text += fmt.Sprintf("[Website](%s) ", *tokenInfo.Website)
	} else {
		text += "(website)"
	}

	var explorerFromInfo string
	if tokenInfo.Explorer != nil {
		explorerFromInfo = *tokenInfo.Explorer
	}

	explorerFromID, err := coin.GetCoinExploreURL(chain, tokenID, tokenType)
	if err != nil {
		return fmt.Sprintf("Failed to retrieve explore url: %v", err)
	}

	if strings.EqualFold(explorerFromInfo, explorerFromID) {
		text += fmt.Sprintf("[Explorer](%s)", explorerFromID)
	} else {
		text += fmt.Sprintf("[ExplorerFromInfo](%s) [ExplorerFromID](%s)", explorerFromInfo, explorerFromID)
	}

	if tokenInfo.Tags != nil {
		text += fmt.Sprintf("\nTags: %s", strings.Join(tokenInfo.Tags, ", "))
	}

	msg := e.checkLogo(logoURL)

	msg += e.checkAssetInfo(tokenInfo)
	if msg == "" {
		text += "\n✅ Check OK"
	} else {
		text += fmt.Sprintf("\nToken check error: \n%s\n", msg)
	}

	text += fmt.Sprintf("\n\n%s\n\n", getLogoHTML(logoURL))

	return text
}

func (e Handler) checkAssetInfo(tokenInfo *assetsmanager.AssetValidationReq) string {
	result, err := e.assetsManager.ValidateAssetInfo(tokenInfo)
	if err != nil {
		log.Debugf(err.Error())

		return "failed to check asset info.json"
	}

	if result.Status == "ok" {
		return ""
	}

	var text string
	for _, assetError := range result.Errors {
		text += fmt.Sprintf("❌ %s\n", assetError.Message)
	}

	return text
}

func (e Handler) checkLogo(url string) string {
	data, err := http.GetHTTPResponseBytes(url)
	if err != nil {
		return fmt.Sprintf("failed to get logo: %s\n", err.Error())
	}

	var text string

	w, h, err := image.GetPNGImageDimensionsFromReader(bytes.NewReader(data))
	if err != nil {
		return fmt.Sprintf("failed to get logo dimensions: %s\n", err.Error())
	}

	err = validation.ValidateImageDimension(w, h)
	if err != nil {
		text += fmt.Sprintf("❌ %s\n", err.Error())
	}

	err = validation.ValidateLogoStreamSize(data)
	if err != nil {
		text += fmt.Sprintf("❌ %s\n", err.Error())
	}

	return text
}
