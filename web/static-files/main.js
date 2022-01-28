const appName = "TW Assets Management";
const appNameHtml = `<a href="https://assets.trustwallet.com">${appName}</a>`;
const gitHub = "https://github.com";
const mainRepoOwner = "trustwallet";
const mainRepoName = "assets";
const mainRepoFullName = mainRepoOwner + "/" + mainRepoName;
const mainRepoUrl = `${gitHub}/${mainRepoFullName}.git`;
const prBodyFooter = `\n\nPR created by ${appNameHtml}`;

const sampleEthContract = "0x6d84682C82526E245f50975190EF0Fff4E4fc077";
const testLogoUrls = [
    "https://raw.githubusercontent.com/trustwallet/assets/master/blockchains/ethereum/assets/0x47F32f9eBFc49a1434eB6190d5D8a80A2Dc36af5/logo.png",
    "https://raw.githubusercontent.com/trustwallet/assets/master/blockchains/ethereum/assets/0x9B9087756eCa997C5D595C840263001c9a26646D/logo.png",
    "https://raw.githubusercontent.com/trustwallet/assets/master/blockchains/ethereum/assets/0xd4c435F5B09F855C3317c8524Cb1F586E42795fa/logo.png",
    "https://raw.githubusercontent.com/trustwallet/assets/master/blockchains/ethereum/assets/0x86876a5fCAcb52a197f194A2c8b2166Af327a6da/logo.png",
    "https://raw.githubusercontent.com/trustwallet/assets/master/blockchains/ethereum/assets/0xD5525D397898e5502075Ea5E830d8914f6F0affe/logo.png"
];
const assetsAPI = "https://api.assets.trustwallet.com/v1/github/oauth"

function addLog(message) {
    console.log(message);
    //const newvalue = document.getElementById("log").value + message + "\n"
    //document.getElementById("log").value = newvalue;
}

async function myAlert(message) {
    addLog(message);
    alert(message);
}

function authHeaders(userToken) {
    return { authorization: `token ${userToken}` };
}

function getQueryParam(param) {
    var urlParams = new URLSearchParams(window.location.search);
    return urlParams.get(param);
}

function getAppMode() {
    if (getQueryParam('maintainer') ||
        window.location.pathname.startsWith('/maintainer') ||
        window.location.pathname.startsWith('/full')
    ) { return 'maintainer'; }
    if (window.location.pathname.startsWith('/list') ||
        window.location.pathname.startsWith('/search')
    ) { return 'search'; }
    return 'contributor';
}

function getUserToken() {
    const tokenQP = getQueryParam("ghtoken");
    const tokenLS = localStorage.getItem("ghtoken");
    if (!tokenQP) {
        if (tokenLS) {
            console.log("Using GH user token from local storage", tokenLS);
            return tokenLS;
        }
        return null;
    }
    // not in LS, but in QP, use it, save it to LS, remove from QP
    localStorage.setItem("ghtoken", tokenQP);
    console.log("Saved GH user token to local storage", tokenQP);
    window.location.search = updateQueryParams('ghtoken', '');
    return tokenQP;
}

function resetUserToken() {
    localStorage.setItem('ghtoken', '');
    window.location.search = updateQueryParams('ghtoken', '');
    console.log("GH user token cleared from local storage and query parameters");
}

function updateQueryParams(newParamName, newValue) {
    var urlParams = new URLSearchParams(window.location.search);
    if (newValue) {
        urlParams.set(newParamName, newValue);
    } else {
        urlParams.delete(newParamName, newValue);
    }
    return urlParams.toString();
}

async function getUser(userToken) {
    const result = await script.octoRequest("GET /user", { headers: authHeaders(userToken) });
    return result.data;
}

async function checkUser(userToken) {
    if (!userToken) {
        return "";
    }
    const user = await getUser(userToken);
    if (!user || !user.login) {
        return "";
    }
    return user.login;
}


async function getUserRepos(userToken, page) {
    const result = await script.octoRequest("GET /user/repos", {
        headers: authHeaders(userToken),
        page: page,
        per_page: 100,
    });
    //console.log(result);
    return result.data;
}

async function getRepo(userToken, owner, repo) {
    const result = await script.octoRequest("GET /repos/:owner/:repo", {
        headers: authHeaders(userToken),
        owner: owner,
        repo: repo
    });
    return result.data;
}

async function isRepoForkOfAssets(repoName, parentRepoFullName, userToken, loginName) {
    try {
        const r2 = await getRepo(userToken, loginName, repoName);
        if (r2 && r2.parent && r2.parent.full_name === parentRepoFullName) {
            return true;
        }
    } catch (error) {
        addLog(`Error: ${repoName} ${error}`);
    }
    return false;
}

async function checkRepo(userToken, loginName, progressCallback) {
    if (!userToken || !loginName) {
        return null;
    }

    // try repo named 'assets' directly
    if (await isRepoForkOfAssets(mainRepoName, mainRepoFullName, userToken, loginName)) {
        addLog(`Fork repo found fast, ${loginName} ${mainRepoName}`);
        return mainRepoName;
    }

    // enumerate all repos of user
    var page = 1;
    var reposGotCount = 0;
    var reposCheckedCount = 0;
    while (page < 10) {
        const repos = await getUserRepos(userToken, page);
        if (!repos || repos.length == 0) {
            break;
        }
        // try all repos
        reposGotCount += repos.length;
        for (const r of repos) {
            reposCheckedCount += 1;
            try {
                if (progressCallback) { progressCallback(r.name); }
                if (r && r.fork && r.name) {
                    if (await isRepoForkOfAssets(r.name, mainRepoFullName, userToken, loginName)) {
                        addLog(`Fork repo found, ${loginName} ${r.name} (checked ${reposCheckedCount}/${reposGotCount})`);
                        return r.name;
                    }
                }
            } catch (error) {
                addLog(`Error: ${r.name} ${error}`);
            }
        }
        page += 1;
    }
    if (reposGotCount == 0) {
        addLog(`Warning: No repositories found for user ${loginName}`);
    } else {
        addLog(`Warning: Fork repo not found, checked ${reposGotCount} repos of user ${loginName}`);
    }
    return null;
}

// Retrieve the size of an image (url or stream) but placing in a hidden img and retrieving size
async function getImageDimension(url, stream = null) {
    try {
        const img = document.getElementById("image-placeholder-for-size-computation");
        //img.src = "";
        if (stream) {
            img.src = "data:image/gif;base64," + stream;
        } else {
            img.src = url;
        }
        return { x: img.naturalWidth, y: img.naturalHeight };
    } catch (error) {
        console.log("getImageDimension exception", error);
        return { x: 0, y: 0 }
    }
}

// return [ArrayBuffer, content-type]
async function logoStreamFromUrl(testLogoUrl) {
    if (!testLogoUrl) {
        return [null, null];
    }
    const response = await fetch(testLogoUrl);
    if (!response || !response.status == 200) {
        addLog(`Logo stream error ${url}`);
        return [null, null];
    }
    return [await response.arrayBuffer(), response.headers.get('Content-Type')];
}

function arrayBufferToBase64(buffer) {
    var binary = '';
    var bytes = new Uint8Array(buffer);
    var len = bytes.byteLength;
    for (var i = 0; i < len; i++) {
        binary += String.fromCharCode(bytes[i]);
    }
    return window.btoa(binary);
}

async function createFork(userToken, owner, repo) {
    const result = await script.octoRequest("POST /repos/:owner/:repo/forks", {
        headers: authHeaders(userToken),
        owner: owner,
        repo: repo
    });
    return result.data.full_name;
}

async function doCreateFork(userToken) {
    const res = await createFork(userToken, mainRepoOwner, mainRepoName);
    myAlert(`Fork created: ${res}`);
    // wait a bit before reload
    new Promise(resolve => setTimeout(resolve, 3000));
    // reload
    window.location = window.location;
}

function newBranchName() {
    const today = new Date();
    const date = (today.getMonth() + 1).toString().padStart(2, '0') + today.getDate().toString().padStart(2, '0') +
        '-' + today.getHours().toString().padStart(2, '0') + today.getMinutes().toString().padStart(2, '0') + today.getSeconds().toString().padStart(2, '0');
    return "br" + date;
}

async function getBranchRef(userToken, loginName, repo, branch) {
    const result = await script.octoRequest("GET /repos/{owner}/{repo}/git/ref/{ref}", {
        headers: authHeaders(userToken),
        owner: loginName,
        repo: repo,
        ref: "heads/" + branch
    });
    return result.data.object.sha;
}

async function createBranch(userToken, loginName, repo, branchName) {
    const ref = await getBranchRef(userToken, loginName, repo, "master");
    const result = await script.octoRequest("POST /repos/:owner/:repo/git/refs", {
        headers: authHeaders(userToken),
        owner: loginName,
        repo: repo,
        ref: `refs/heads/${branchName}`,
        sha: ref
    });
    return result.data.ref;
}

async function createBlob(userToken, loginName, repo, content, encoding) {
    const result = await script.octoRequest("POST /repos/:owner/:repo/git/blobs", {
        headers: authHeaders(userToken),
        owner: loginName,
        repo: repo,
        content: content,
        encoding: encoding
    });
    return result.data.sha;
}

async function createFiles(userToken, loginName, repo, basePath, subFolder, logoContents, infoJsonContent) {
    var fileInfos = [];

    const folder = basePath + "/" + subFolder;

    const logoSha = await createBlob(userToken, loginName, repo, logoContents, "base64");
    if (!logoSha) {
        return null;
    }
    fileInfos.push({ path: folder + "/logo.png", sha: logoSha, mode: "100644", type: "blob" });

    fileInfos.push({ path: folder + "/info.json", content: infoJsonContent, mode: "100644", type: "blob" });

    return fileInfos;
}

async function createTree(userToken, loginName, repo, baseSha, fileInfos) {
    const result = await script.octoRequest("POST /repos/:owner/:repo/git/trees", {
        headers: authHeaders(userToken),
        owner: loginName,
        repo: repo,
        tree: fileInfos,
        base_tree: baseSha
    });
    return result.data.sha;
}

async function createCommit(userToken, loginName, repo, baseSha, tree, tokenName) {
    const result = await script.octoRequest("POST /repos/:owner/:repo/git/commits", {
        headers: authHeaders(userToken),
        owner: loginName,
        repo: repo,
        message: `Info for token ${tokenName}`,
        tree: tree,
        parents: [baseSha]
    });
    return result.data.sha;
}

async function updateReference(userToken, loginName, repo, ref, sha) {
    const result = await script.octoRequest("POST /repos/:owner/:repo/git/refs/:ref", {
        headers: authHeaders(userToken),
        owner: loginName,
        repo: repo,
        ref: ref,
        sha: sha
    });
    return result.data.object.sha;
}

async function createPull(userToken, loginName, repo, branchName, tokenName, tokenType, debugPrTargetFork) {
    const title = `Add ${tokenName} (${tokenType})`;
    const text = `Adding info for ${tokenType} token '${tokenName}'.\n\n${prBodyFooter}`;
    var owner1 = mainRepoOwner;
    var repo1 = mainRepoName;
    if (debugPrTargetFork) {
        owner1 = loginName;
        repo1 = repo;
    }
    const result = await script.octoRequest("POST /repos/:owner/:repo/pulls", {
        headers: authHeaders(userToken),
        owner: owner1,
        repo: repo1,
        title: title,
        head: loginName + ":" + branchName,
        base: "master",
        body: text,
        maintainer_can_modify: true,
        draft: false
    });
    return result.data.number;
}

async function getPulls(userToken, owner, repo) {
    const result = await script.octoRequest("GET /repos/:owner/:repo/pulls", {
        headers: authHeaders(userToken),
        owner: owner,
        repo: repo,
        per_page: 100,
    });
    return result.data;
}

async function getPrFiles(userToken, prNum) {
    const url = `https://api.github.com/repos/${mainRepoFullName}/pulls/${prNum}/files`;
    let resp = await fetch(url, { headers: authHeaders(userToken) });
    if (resp.status != 200) {
        myAlert(`Error from ${url}, status ${resp.status} ${resp.statusText}`);
        return [];
    }
    const respJson = await resp.json();
    const files = respJson.map(e => e.filename);
    return files;
}

async function checkUrlByBackend(url) {
    const beUrl = `/checkUrl?url=${encodeURI(url)}`;
    console.log(`checkUrlByBackend ${beUrl}`);
    let resp = await fetch(beUrl);
    return resp.status;
}

function start() {
    // Preview for a single logo, supports url and stream, rounded mask, dimming
    Vue.component('logo-single-preview', {
        props: {
            size: Number,
            logourl: String,
            logostream: String,
            dimmed: Boolean,
            rounded: Boolean
        },
        computed: {
            src: function () {
                if (this.logostream) {
                    // image data inline
                    return "data:image/gif;base64," + this.logostream;
                }
                if (this.logourl) { return this.logourl; }
                // fallback TODO replace with placeholder image
                return "img/emptylogo.png";
            },
            opacity: function () { return this.dimmed ? 0.6 : 1; },
            roundedRadius: function () { return this.rounded ? "48%" : "0%" },
            style: function () { return `max-width: 64; opacity: ${this.opacity}; border-radius: ${this.roundedRadius};`; }
        },
        data: function () {
            return {}
        },
        template: `<img :width="size" :height="size" :src="src" :style="style" />`
    });

    // Token text for preview
    Vue.component('logo-preview-with-name', {
        props: {
            logourl: String,
            logostream: String,
            name: String,
            dimmed: Boolean,
        },
        template:
            `
                <tr>
                    <td style="padding: 5;">
                        <logo-single-preview :size="32" :logourl="logourl" :logostream="logostream" :dimmed="dimmed" :rounded="true" />
                    </td>
                    <td class="wide" style="align: right;">
                        <span style="vertical-align: center; font-size: 70%;">{{name ? name.substring(0, 18) : '(Token)'}}</span>
                    </td>
                </tr>
            `
    });

    // Logo preview  with a few other logos above/below
    Vue.component('logo-column-preview', {
        props: {
            logourl: String,
            logostream: String,
            tokenname: String,
            textcolor: String,
        },
        data: function () {
            return {
                logosBuiltin: {
                    Bitcoin: "https://raw.githubusercontent.com/trustwallet/assets/master/blockchains/bitcoin/info/logo.png",
                    Trust: "https://raw.githubusercontent.com/trustwallet/assets/master/blockchains/binance/assets/TWT-8C2/logo.png",
                    Binance: "https://raw.githubusercontent.com/trustwallet/assets/master/blockchains/binance/info/logo.png",
                    Ethereum: "https://raw.githubusercontent.com/trustwallet/assets/master/blockchains/ethereum/info/logo.png"
                }
            }
        },
        template:
            `
                <div style="display: inline-block; padding: 15px; align: right;">
                    <table :style="'color: ' + textcolor + '; width: 130px;'">
                        <logo-preview-with-name :logourl="logosBuiltin['Bitcoin']" name="Bitcoin" :dimmed="true" />
                        <!-- logo-preview-with-name :logourl="logosBuiltin['Trust']" name="Trust" :dimmed="true" / -->
                        <logo-preview-with-name :logourl="logourl" :logostream="logostream" :name="tokenname" :dimmed="false" />
                        <!-- logo-preview-with-name :logourl="logosBuiltin['Binance']" name="Binance" :dimmed="true" / -->
                        <logo-preview-with-name :logourl="logosBuiltin['Ethereum']" name="Ethereum" :dimmed="true" />
                    </table>
                </div>
            `
    });

    // Log preview with a column with light background and one with dark background
    Vue.component('logo-preview', {
        props: {
            logourl: String,
            logostream: String,
            tokenname: String,
        },
        template:
            `
                <div>
                    <table class="left">
                        <tr>
                            <td style="background-color: #E0E0E0;">
                                <logo-column-preview :logourl="logourl" :logostream="logostream" :tokenname="tokenname" textcolor="black" />
                            </td>
                            <td style="background-color: #202020; color: white;">
                                <logo-column-preview :logourl="logourl" :logostream="logostream" :tokenname="tokenname" textcolor="white" />
                            </td>
                        </tr>
                    </table>
                </div>
            `
    });

    Vue.component('token-view', {
        props: {
            tokenInfo: Object
        },
        data: function () {
            return {
                circulationHolders: '',
                script: script,
            }
        },
        updated: async function () {
            this.circulationHolders = await script.assets.getTokenCirculationSafe(this.tokenInfo.type, this.tokenInfo.contract, true);
        },
        computed: {
            links: function () {
                if (!this.tokenInfo) { return ""; }
                let links = "";
                if (this.tokenInfo.logoUrl) {
                    links += `<a href="${this.tokenInfo.logoUrl}" target="_blank">Logo</a> `;
                } else {
                    links += "(logo) ";
                }
                if (this.tokenInfo.infoUrl) {
                    links += `<a href="${this.tokenInfo.infoUrl}" target="_blank">Info</a> `;
                } else {
                    links += "(info) ";
                }
                const explorer = this.tokenInfo.explorerUrl();
                if (explorer) {
                    links += `<a href="${explorer}" target="_blank">Explorer</a> `;
                } else {
                    links += "(explorer) ";
                }
                const website = this.tokenInfo.info["website"];
                if (website) {
                    links += `<a href="${website}" target="_blank">Website</a> `;
                } else {
                    links += "(website) ";
                }
                return links;
            }
        },
        methods: {
            checkInfoButton: async function () {
                addLog("starting TokenInfo check...");
                let [resnum, resmsg] = await script.assets.checkTokenInfo(this.tokenInfo,
                    { checkUrl: checkUrlByBackend },
                    { get: getImageDimension }, true);
                if (resnum >= 2) {
                    myAlert("Check result: ERROR \n" + resmsg);
                } else if (resnum >= 1) {
                    myAlert("Check result: Warning \n" + resmsg);
                } else {
                    myAlert("Check OK: \n" + resmsg);
                }
            },
        },
        template: `
            <div class="col-12 col-lg-7 mt-2">
                    <h4 class="font-weight-bold">Token preview</h4>
                    <div class="mt-2 mb-2"><logo-preview :logourl="tokenInfo.logoUrl" :logostream="tokenInfo.logoStream"/></div>
                    <div class="table-responsive-sm">
                <table class="table table-borderless">
                    <tr>
                        <td>Type:</td>
                        <td><strong>{{tokenInfo.type}}</strong></td>
                    </tr>
                    <tr>
                        <td>Contract:</td>
                        <td><strong>{{tokenInfo.contract}}</strong></td>
                    </tr>
                    <tr v-show="tokenInfo.assetId">
                        <td>Asset ID:</td>
                        <td>
                            <span class="smallfont" v-text="tokenInfo.assetId"/>
                        </td>
                    </tr>
                    <tr>
                        <td colspan="2">
                        <div class="form p-0 mt-2 mb-2">
  <textarea class="form-control" placeholder="" id="info.info" readonly="true" rows="8">{{tokenInfo.infoString}}</textarea>
                      </div>

                        </td>
                    </tr>
                    <tr>
                        <td colspan="2">
                            Links: <span v-html="links"></span>
                        </td>
                    </tr>
                    <tr>
                        <td>Circulation (holders):</td>
                        <td>{{circulationHolders}}</td>
                    </tr>

                              </div>
                </table>
                    <button class="btn btn-outline-primary" type="button" v-on:click="checkInfoButton">Check</button>
            </div>
        `
    });

    Vue.component('pull-label-item', {
        props: {
            label: Object,
        },
        computed: {
            initials: function () {
                const words = this.label.name.split(' ');
                inits = words.map(w => w[0]).join('');
                return inits.substring(0, 3);
            },
        },
        template:
            `
                <span class="pull-label-small" v-bind:style="{ backgroundColor: label.color }"
                    v-bind:title="label.name">{{initials}}
                </span>
            `
    });

    Vue.component('pull-item', {
        props: {
            pull: Object,
        },
        methods: {
            onclick: function () {
                this.$emit('select', this.pull);
            },
        },
        computed: {
            shortTitle: function () {
                const maxlen = 25;
                if (this.pull.title.length <= maxlen) {
                    return this.pull.title;
                }
                return this.pull.title.substring(0, maxlen - 1) + '...';
            }
        },
        template:
            `
                <div class="pull-item selection-item" v-on:click="onclick">
                    <strong>{{pull.number}}</strong>
                    <span v-bind:title="pull.title">{{shortTitle}}</span>
                    <pull-label-item v-for="label in pull.labels" :label="label"></pull-label-item>
                </div>
            `
    });

    Vue.component('pull-content', {
        props: {
            pull: Object,
            tokenIds: Array,
        },
        data: function () {
            return {
                selectedToken: Object,
                mainRepoFullName: mainRepoFullName,
                gitHub: gitHub,
            }
        },
        async updated() {
            //console.log("pull-content updated", this.tokenIds.length);
            if (!this.tokenIds || this.tokenIds.length == 0) {
                if (this.selectedToken) {
                    this.selectedToken = null;
                    this.$emit('select', this.selectedToken);
                }
            } else if (this.tokenIds && this.tokenIds.length >= 1) {
                if (!this.selectedToken) {
                    this.selectedToken = this.tokenIds[0];
                    this.$emit('select', this.selectedToken);
                }
            }
        },
        methods: {
            onclick: function (e) {
                const token = JSON.parse(e.target.attributes.token.value);
                //console.log("token", token);
                this.selectedToken = token;
                //console.log("selectedToken", this.selectedToken.id, this.selectedToken);
                this.$emit('select', token);
            },
        },
        template:
            `
                <div>
                    <div class="margined">
                        <div class="margined"><a :href="gitHub + '/' + mainRepoFullName + '/pull/' + pull.number" target="_blank" rel="noopener noreferrer"><strong>PR #{{pull.number}}</strong></a>: {{pull.title}}</div>
                        <div class="margined">repo: <a :href="gitHub + '/'+ pull.owner + '/' + pull.repo" target="_blank" rel="noopener noreferrer">{{pull.owner}}/{{pull.repo}}</a>/{{pull.branch}}</div>
                        <div class="margined">Detected token(s) in PR:</div>
                    </div>
                    <div id="pull-content-list" class="selection-list margined ">
                        <div v-for="token in tokenIds" class="selection-item" v-on:click="onclick" :token="JSON.stringify(token)" :tokentype="token.type" :tokenid="token.id"
                            v-bind:class="{'selected-item': token && selectedToken && token.id == selectedToken.id}">{{token.type}}: {{token.id}}</div>
                    </div>
                </div>
            `
    });

    Vue.component('pulls-list', {
        props: {
            userToken: String,
            enabled: Boolean,
        },
        data: function () {
            return {
                pulls: Array,
                selected: Object,
                selectedPullTokenIds: Array,
                script: script,
            }
        },
        async created() {
            // enabled is set delayed, check queryparam here directly
            appMode = getAppMode();
            if (this.enabled || appMode === 'maintainer') {
                await this.loadOpenPrs();
            }
        },
        methods: {
            loadOpenPrs: async function () {
                if (!this.userToken) { this.pulls = []; return; }
                const rawPulls = await getPulls(this.userToken, mainRepoOwner, mainRepoName);
                this.pulls = rawPulls.map(p => ({
                    number: p.number,
                    title: p.title,
                    owner: p.head.user ? p.head.user.login : '',
                    repo: p.head.repo ? p.head.repo.name : '',
                    branch: p.head.ref,
                    labels: p.labels.map(l => ({ name: l.name, color: l.color })),
                }));
                if (this.pulls && this.pulls.length >= 1) {
                    await this.onselectpull(this.pulls[0], true);
                }
            },
            onselectpull: async function (pr, silentMode = false) {
                this.selected = pr;
                //console.log(`select ${pr.number}: ${pr.owner}/${pr.repo}/${pr.branch}`);
                this.selectedPullTokenIds = [];
                this.$emit('selecttoken', new script.assets.TokenInfo());

                const files = await getPrFiles(this.userToken, pr.number);
                if (!files || files.length == 0) {
                    myAlert(`Could not retrieve files from PR ${pr.number}`);
                    return;
                }
                const ids = script.assets.tokenIdsFromFiles(files);
                if (!ids || ids.length == 0) {
                    if (!silentMode) {
                        myAlert(`Could not retrieve tokens from PR ${pr.number}`);
                    }
                    return;
                }
                this.selectedPullTokenIds = ids.map(i => ({ type: i[0], id: i[1], owner: pr.owner, repo: pr.repo, branch: pr.branch }));
            },
            onselecttoken: function (token) {
                this.$emit('selecttoken', token);
            },
        },
        template:
            `
                <div>
                    <div class="margined">
                        <button v-on:click="loadOpenPrs()">Search</button>
                        <span>{{pulls.length}} PRs</span>
                    </div>
                    <div id="pulls-list-list" class="selection-list margined">
                        <pull-item v-for="pull in pulls" :pull="pull" v-on:select="onselectpull"
                            v-bind:class="{'selected-item': pull.number == selected.number}"></pull-item>
                    </div>
                    <pull-content :pull="selected" :token-ids="selectedPullTokenIds" v-on:select="onselecttoken"></pull-content>
                </div>
            `
    });

    Vue.component('token-search-item', {
        props: {
            token: Object,
        },
        methods: {
            onclick: function () {
                this.$emit('select', this.token);
                //console.log("token", this.token.token_id, this.token);
            }
        },
        computed: {
            symbol: function () { return this.token.symbol.substring(0, 10); },
            name: function () { return this.token.name.substring(0, 18); },
            type: function () { return this.token.type.toLowerCase(); },
        },
        template:
            `
                <div class="token-search-item selection-item" v-on:click="onclick">
                    <strong>{{symbol}}</strong>
                    {{name}}
                    ({{type}})
                </div>
            `
    });

    Vue.component('token-search', {
        data: function () {
            return {
                queryString: '',
                //includeUnverified: false,
                tokens: Array,
                selected: Object,
            }
        },
        methods: {
            searchButton: async function () {
                if (this.queryString.length < 2) {
                    // too short, do not search
                    return;
                }
                this.tokens = [];
                //const coinTypes = "60,195,714,20000714"; // eth, tron, binance, bep20; networks=${coinTypes}
                //${this.includeUnverified ? '&include_unverified=true' : ''}
                const url = `https://api.trustwallet.com/v1/search/assets?query=${this.queryString}&version=100&type=all`;
                let resp = await fetch(url);
                if (resp.status != 200) {
                    myAlert(`Error from ${url}, status ${resp.status} ${resp.statusText}`);
                    return;
                }
                const respJson = await resp.json();
                var h = '';
                if (respJson && respJson.docs) {
                    this.tokens = respJson.docs.map(t => t);
                }
                if (this.tokens && this.tokens.length >= 1) {
                    await this.onselecttoken(this.tokens[0]);
                }
            },
            onselecttoken: async function (token) {
                this.selected = token;
                //console.log("selected", this.selected, this.selected.type, this.selected.type.token_id);
                this.$emit('selecttoken', this.selected);
            },
        },
        template:
            `
                <div id="token-search">
                <div class="row align-items-center pb-2">
                        <div class="d-inline-flex"><input id="search-input" class="form-control m-1" placeholder="Token" v-on:change="searchButton" v-model="queryString">
                        <button type="submit" class="btn btn-primary m-1" v-on:click="searchButton">Search</button></div>
                </div>
                    <div id="token-search-list" class="c">
                        <token-search-item v-for="token in tokens" :token="token" v-on:select="onselecttoken"
                            v-bind:class="{'selected-item': token && selected && token.address == selected.address}"></token-search-item>
                    </div>
                </div>
            `
    });

    Vue.component('link-item-edit', {
        props: {
            tokenInput: script.assets.TokenInput,
            name: String,
            prefix: String,
        },
        data: function () {
            return {
                value: ''
            }
        },
        computed: {
            link: function () {
                for (i in this.tokenInput.links) {
                    if (this.tokenInput.links[i].name === this.name) {
                        // fill with default if empty
                        if (this.prefix && !this.tokenInput.links[i]['url']) {
                            this.tokenInput.links[i]['url'] = this.prefix;
                        }
                        return this.tokenInput.links[i];
                    }
                }
                return null;
            },
        },
        template:
            `
                <tr>
                    <td>{{name}}:</td>
                    <td>
                        <input v-model="link['url']" class="input wide" :placeholder="prefix"/>
                    </td>
                </tr>
            `
    });

    Vue.component('links-add', {
        props: {
            tokenInput: script.assets.TokenInput
        },
        data: function () {
            return {
                selectedType: ''
            }
        },
        computed: {
            linksText: function () {
                if (this.tokenInput.links.length == 0) {
                    return '(no links)';
                }
                let t = '';
                for (i in this.tokenInput.links) {
                    t += this.tokenInput.links[i].url + ' ';
                }
                return t;
            },
            // link types that are not added
            linkTypes: function () {
                return Object.keys(script.assets.LinksKeys).filter(l => !this.linkAdded(l));
            },
            linkPrefixes: function () {
                return script.assets.LinksKeys;
            }
        },
        async created() {
            this.selectedType = this.linkTypes[0];
        },
        methods: {
            linkAdded: function (name) {
                return this.tokenInput.links.find(x => x.name === name) !== undefined;
            },
            addLink: function () {
                const type = this.selectedType;
                if (type) {
                    const prefix = this.linkPrefixes[type];
                    const linkItem = {
                        name: type,
                        url: null,
                        __prefix: prefix,
                    };
                    this.tokenInput.addLinkItem(linkItem);
                    this.selectedType = this.linkTypes[0];
                }
            }
        },
        template:
            `
                <tr>
                    <td>Links:</td>
                    <td class="d-flex align-items-center">
                        <select class="input" v-model="selectedType">
                            <option v-for="linkType in linkTypes">{{linkType}}</option>
                        </select>
                        <button class="btn btn-outline-primary btn-sm" type="button" v-on:click="addLink()">Add Link</button>
                        <span class="smallfont m-2">Add all that apply</span>
                    </td>
                </tr>
            `
    });

    Vue.component('tags-add', {
        props: {
            tokenInput: script.assets.TokenInput
        },
        data: function () {
            return {
                tagsText: '',
                selectedTag: '', // tag selected in the dropdown
            }
        },
        async created() {
            this.tagsText = this.arrayToText(this.tokenInput.tags);
            this.selectedTag = '';
        },
        methods: {
            updateFromText: function () {
                this.tokenInput.tags = this.textToArray(this.tagsText);
                this.updateText();
            },
            updateText: function () {
                this.tagsText = this.arrayToText(this.tokenInput.tags);
            },
            toggleTag: function () {
                if (!this.selectedTag) {
                    return;
                }
                if (this.tokenInput.tags.includes(this.selectedTag)) {
                    // remove
                    this.tokenInput.tags = this.tokenInput.tags.filter(t => (t !== this.selectedTag));
                } else {
                    // add
                    this.tokenInput.tags.push(this.selectedTag);
                }
                this.updateText();
            },
            arrayToText: function (tags) {
                let t = '';
                for (i in tags) {
                    if (t) { t += ' '; }
                    t += tags[i];
                }
                return t;
            },
            textToArray: function(text) {
                text = text.replace(',', '');
                let arr = text.split(' ');
                return arr.filter(t => (t.length > 0)); // filter empty ones
            }
        },
        template:
            `
                <tr>
                    <td>Tags:</td>
                    <td>
                        <select class="input" v-model="selectedTag" v-on:click="toggleTag()">
                            <option value="">(add/remove a tag by selecting from the list)</option>
                            <option v-for="tag in script.assets.TagValues" :value="tag.id">{{tag.id}} <!-- {{tag.description}} --></option>
                        </select>
                        <input v-model="tagsText" class="input" v-on:change="updateFromText()" placeholder="Select at least one tag" />
                    </td>
                </tr>
            `
    });

    Vue.component('main-add-token', {
        props: {
            userToken: String,
            loginName: String,
            repo: String,
        },
        data: function () {
            return {
                script: script,
                tokenInput: new script.assets.TokenInput(),
                errorLogo: ' ',
                errorContract: ' ',
                errorWebsite: ' ',
                errorExplorer: ' ',
                fixedContract: null,
                fixedWebsite: null,
                fixedExplorer: null,
                inputLogoFilename: '',
                inputLogoDetails: '',
                debugMode: false,
                debugPrTargetFork: false,
                testLogoIndex: 0,
                tokenInfo: new script.assets.TokenInfo(),
                checkButtonText: '',
                prButtonText: '',
                mainRepoUrl: mainRepoUrl,
                gitHub: gitHub,
            }
        },
        async created() {
            this.clearInput();
            this.debugMode = getQueryParam("debug");
        },
        methods: {
            clearInput: function () {
                this.tokenInput = new script.assets.TokenInput();
                this.inputLogoSetStream(null, null, 0, null);
            },
            isInputStillEmpty: function () {
                // don't display check result if all inputs are empty
                if (!this.tokenInput.name && !this.tokenInput.contract && !this.tokenInput.logoStream) {
                    return true;
                }
                return false;
            },
            checkInputButton: async function () {
                this.checkButtonText = 'Checking ...';
                addLog("starting TokenInput check...");
                const [resnum, resmsg, fixed] = await script.assets.checkTokenInput(this.tokenInput,
                    { checkUrl: checkUrlByBackend }, { get: getImageDimension }, true);
                if (resnum == 0) {
                    myAlert("Check result: OK\n" + resmsg);
                    this.checkButtonText = '';
                    return resnum;
                }
                if (!fixed) {
                    // cannot be fixed
                    if (resnum >= 2) {
                        myAlert("Check result: ERROR\n" + resmsg);
                        this.checkButtonText = '';
                        return resnum;
                    }
                    myAlert("Check result: Warning\n" + resmsg);
                    this.checkButtonText = '';
                    return resnum;
                }
                // can be fixed
                this.tokenInput = fixed;
                myAlert("Check result: Fixed! \n " + resmsg);
                this.checkButtonText = '';
                return resnum;
            },
            tokenInputLogoChanged: async function () {
                if (this.isInputStillEmpty()) {
                    return;
                }
                const [resnum, resmsg] = await script.assets.checkTokenInputLogo(this.tokenInput, { get: getImageDimension });
                this.errorLogo = ' ';
                if (resnum >= 2) {
                    this.errorLogo = resmsg;
                }
            },
            tokenInputContractChanged: async function () {
                if (this.isInputStillEmpty()) {
                    return;
                }
                // auto-fill explorer
                if (this.tokenInput.contract && !this.tokenInput.explorerUrl) {
                    const explorer = script.assets.explorerUrlForToken(this.tokenInput.type, this.tokenInput.contract);
                    this.tokenInput.explorerUrl = explorer;
                }

                const [resnum, resmsg, fixed] = await script.assets.checkTokenInputContract(this.tokenInput);
                this.fixedContract = fixed;
                this.errorContract = ' ';
                if (resnum >= 2 || fixed) {
                    this.errorContract = resmsg;
                } else {
                    // contract OK, try to preload info from it
                    const tokenInf = await script.assets.getExternalTokenInfo(this.tokenInput.type, this.tokenInput.contract);
                    if (tokenInf) {
                        if (tokenInf.name && tokenInf.symbol && !this.tokenInput.name) { this.tokenInput.name = `${tokenInf.name} (${tokenInf.symbol})`; }
                        if (tokenInf.name && !this.tokenInput.name) { this.tokenInput.name = tokenInf.name; }
                        if (tokenInf.symbol && !this.tokenInput.symbol) { this.tokenInput.symbol = tokenInf.symbol; }
                        if (tokenInf.decimals && !this.tokenInput.decimals) {
                            this.tokenInput.decimals = script.assets.safeParseInt(tokenInf.decimals);
                        }
                        if (tokenInf.website && !this.tokenInput.website) { this.tokenInput.website = tokenInf.website; }
                    }
                }
            },
            tokenInputNameChanged: async function () {
            },
            tokenInputSymbolChanged: async function () {
            },
            tokenInputDecimalsChanged: async function () {
                // enforce numerical
                this.tokenInput.decimals = script.assets.safeParseInt(this.tokenInput.decimals);
                this.tokenInput.decimals = Math.max(0, this.tokenInput.decimals);
                this.tokenInput.decimals = Math.min(25, this.tokenInput.decimals);
            },
            tokenInputWebsiteChanged: async function () {
                if (this.isInputStillEmpty()) {
                    return;
                }
                const [resnum, resmsg, fixed] = await script.assets.checkTokenInputWebsite(this.tokenInput, { checkUrl: checkUrlByBackend });
                this.fixedWebsite = fixed;
                this.errorWebsite = ' ';
                if (resnum >= 2 || fixed) {
                    this.errorWebsite = resmsg;
                }
            },
            tokenInputExplorerChanged: async function () {
                if (this.isInputStillEmpty()) {
                    return;
                }
                const [resnum, resmsg, fixed] = await script.assets.checkTokenInputExplorer(this.tokenInput, { checkUrl: checkUrlByBackend });
                this.fixedExplorer = fixed;
                this.errorExplorer = ' ';
                if (resnum >= 2 || fixed) {
                    this.errorExplorer = resmsg;
                }
            },
            tokenInputDescriptionChanged: async function () {
            },
            tokenInputChanged: async function () {
                await this.tokenInputLogoChanged();
                await this.tokenInputContractChanged();
                await this.tokenInputNameChanged();
                await this.tokenInputSymbolChanged();
                await this.tokenInputDecimalsChanged();
                await this.tokenInputWebsiteChanged();
                await this.tokenInputExplorerChanged();
                await this.tokenInputDescriptionChanged();
            },
            inputLogoSetStream: function (stream, hintName, hintSize, hintMime) {
                if (!stream) {
                    this.tokenInput.logoStream = null;
                    this.tokenInput.logoStreamSize = 0;
                    this.tokenInput.logoStreamType = null;
                    this.inputLogoFilename = "(no image)";
                    this.inputLogoDetails = "";
                    return;
                }
                this.tokenInput.logoStream = stream;

                this.inputLogoFilename = hintName;

                let details = "(";
                this.tokenInput.logoStreamSize = 0;
                if (hintSize && hintSize > 0) {
                    this.tokenInput.logoStreamSize = hintSize;
                    details += `${hintSize} bytes`;
                } else {
                    details += `${this.tokenInput.logoStream.length} base64 bytes`;
                }
                if (hintName) {
                    details += `, ${hintName}`;
                }
                this.tokenInput.logoStreamType = "";
                if (hintMime) {
                    this.tokenInput.logoStreamType = hintMime;
                    details += `, ${hintMime}`;
                }
                details += ")";
                this.inputLogoDetails = details;
            },
            logoFileSelected: async function () {
                const file = document.getElementById("input.file-selector").files[0];
                if (file) {
                    var reader = new FileReader();
                    reader.readAsArrayBuffer(file);
                    app = this;
                    reader.onload = async function (evt) {
                        var contents = evt.target.result;
                        var base64 = arrayBufferToBase64(contents);
                        app.inputLogoSetStream(base64, file.name, file.size, file.type);
                        await app.tokenInputLogoChanged();
                    }
                    reader.onerror = function (evt) {
                        myAlert(`Error reading file ${file.name}`);
                    }
                }
            },
            createPullError: async function (message) {
                myAlert("PR creation error: " + message);
                this.prButtonText = '';
            },
            createBranchAndPull: async function () {
                if (!this.loginName) {
                    this.createPullError("Log in first!");
                    return;
                }
                if (!this.repo) {
                    this.createPullError("Create a fork of the assets repo first!");
                    return;
                }
                if (!this.tokenInput.contract || !this.tokenInput.name || !this.tokenInput.type || !this.tokenInput.symbol || !this.tokenInput.decimals || !this.tokenInput.logoStream) {
                    this.createPullError("Fill in all the fields!");
                    return;
                }
                this.tokenInfo = this.tokenInput.toTokenInfo();
                if (!this.tokenInfo) {
                    this.createPullError("TokenInput error");
                    return;
                }

                this.prButtonText = "creating PR ...";

                const branchName = newBranchName();

                addLog(`name: ${this.tokenInput.name}  type: ${this.tokenInput.type}  contract: ${this.tokenInput.contract}  branch ${branchName}`);

                this.prButtonText = "creating branch ...";
                const branchRef = await createBranch(this.userToken, this.loginName, this.repo, branchName);
                if (!branchRef) {
                    this.createPullError(`Could not create branch ${branchName}`);
                    return;
                }
                addLog(`Created branch ${this.loginName}/${this.repo}/${branchName}`);

                let stream = null;
                if (this.tokenInput.logoStream && this.tokenInput.logoStream.length > 10) {
                    // stream is there, use that
                    stream = this.tokenInput.logoStream;
                }
                if (!stream) {
                    this.createPullError(`Could not retrieve logo contents`);
                    return;
                }

                const chain = script.assets.chainFromAssetType(this.tokenInput.type);
                if (!chain || chain == "unknown") {
                    this.createPullError(`Could not retrieve chain from token type ${this.tokenInput.type}`);
                    return;
                }

                this.prButtonText = "creating files ...";
                const fileInfos = await createFiles(this.userToken, this.loginName, this.repo, `blockchains/${chain}/assets`, this.tokenInput.contract, stream, this.tokenInfo.infoString);
                if (!fileInfos || fileInfos.length == 0) {
                    this.createPullError(`Could not create files`);
                    return;
                }
                addLog(`Created ${fileInfos.length} new files`);

                const branchSha = await getBranchRef(this.userToken, this.loginName, this.repo, branchName);
                if (!branchSha) {
                    this.createPullError(`Could not get ref for branch ${branchName}`);
                    return;
                }

                const tree = await createTree(this.userToken, this.loginName, this.repo, branchSha, fileInfos);
                if (!tree) {
                    this.createPullError(`Could not create tree with files`);
                    return;
                }

                this.prButtonText = "creating commit ...";
                const commit = await createCommit(this.userToken, this.loginName, this.repo, branchSha, tree, this.tokenInput.name);
                if (!commit) {
                    this.createPullError(`Could not create commit`);
                    return;
                }
                addLog(`Created new commit ${commit}`);

                const newBranchSha = await updateReference(this.userToken, this.loginName, this.repo, "heads/" + branchName, commit);
                if (!newBranchSha) {
                    this.createPullError(`Could not update branch ${branchName} to commit ${commit}`);
                    return;
                }

                this.prButtonText = "creating pull request ...";
                const pullNumber = await createPull(this.userToken, this.loginName, this.repo, branchName, this.tokenInput.name, this.tokenInput.type, this.debugPrTargetFork);  // true for debug
                if (!pullNumber) {
                    this.createPullError(`Could not create PR`);
                    return;
                }
                this.prButtonText = '';

                var pullOwner = mainRepoOwner;
                var pullRepo = mainRepoName;
                if (this.debugPrTargetFork) {
                    pullOwner = this.loginName;
                    pullRepo = this.repo;
                }
                const pullUrl = `${gitHub}/${pullOwner}/${pullRepo}/pull/${pullNumber}`;
                myAlert(`Created PR ${pullNumber}   ${pullUrl}`);
            },
            createPullButton: async function () {
                checkResult = await this.checkInputButton();
                let createPR = false;
                if (checkResult == 0) {
                    // OK
                    createPR = true;
                } else {
                    if (confirm(`There were errors/warnings in the input, please fix them first.\nAlternatively, create a PR with errors now?`)) {
                        createPR = true;
                    } else {
                        createPR = false;
                        addLog(`Not creating PR due to check errors, checkResult ${checkResult}`);
                    }
                }
                if (createPR) {
                    await this.createBranchAndPull(this.userToken, this.loginName, this.repo, this.tokenInput, this.debugPrTargetFork);
                }
            },
            debugTestLogoGetNext: async function () {
                this.testLogoIndex = (this.testLogoIndex + 1) % testLogoUrls.length;
                const [streamArray, mime] = await logoStreamFromUrl(testLogoUrls[this.testLogoIndex]);
                const stream = arrayBufferToBase64(streamArray);
                this.inputLogoSetStream(stream, "test" + this.testLogoIndex, streamArray.byteLength, mime);
            },
            debugFillWithDummyData: async function () {
                this.tokenInput.name = "LegumeToken";
                this.tokenInput.type = "erc20";
                this.tokenInput.contract = sampleEthContract;
                this.tokenInput.website = "https://en.wikipedia.org/wiki/Legume";
                this.tokenInput.description = "This is the best-tasting DeFi finance project.";
                this.tokenInput.explorerUrl = `https://etherscan.io/token/${sampleEthContract}`;
                await this.debugTestLogoGetNext();
                await this.tokenInputChanged();
            },
        },
        template:
            `
                <div id="main-add-token" class="row mt-3">
                    <div id="add-input">
                        <div class="">
                            <h4>Add New Token</h4>
                            <p id="add-explanation" class="smallfont">
                                    See guide on <a href="https://developer.trustwallet.com/assets/new-asset"
                                        target="_blank">developer</a> and
                                    <a href="https://community.trustwallet.com/t/how-to-submit-a-token-logo/3863"
                                        target="_blank">community</a> sites.

                                    Fill in the token/project details.
                                    If all is OK, a new
                                    <a href="https://docs.github.com/en/free-pro-team@latest/desktop/contributing-and-collaborating-using-github-desktop/creating-an-issue-or-pull-request"
                                        target="_blank">pull request</a> will be created in your name,
                                    against the <a :href="mainRepoUrl"
                                        target="_blank">assets repo</a>.
                            </p>
                        </div>
                        <div>
                            <div>
                                <logo-preview :logostream="tokenInput.logoStream" v-show="tokenInput.logoStream"
                                    :tokenname="tokenInput.name.substring(0, 16)" />
                            </div>
                            <div class="mt-3">
                                <button class="btn btn-outline-primary"
                                    onclick="document.getElementById('input.file-selector').click();">Upload Logo</button>
                                <input id="input.file-selector" type="file" style="display: none;"
                                    v-on:change="logoFileSelected()" />
                                <span v-html="inputLogoFilename" v-bind:title="inputLogoDetails" id="input.logo-input" class="smallfont"></span>
                            </div>
                            <div class="smallfont error padded">
                                {{errorLogo}}
                            </div>
                            <table class="wide">
                                <tr>
                                    <td>Type:</td>
                                    <td>
                                        <select v-model="tokenInput.type" class="input wide"
                                                    v-on:change="tokenInputChanged()">
                                            <option value="ERC20">ERC20 (Ethereum)</option>
                                            <option value="BEP2">BEP2 (Binance)</option>
                                            <option value="BEP20">BEP20 (Binance Smart Chain)</option>
                                            <option value="TRC10">TRC10 (Tron)</option>
                                            <option value="TRC20">TRC20 (Tron)</option>
                                            <option value="ETC20">ETC20 (Ethereum Classic)</option>
                                            <option value="GO20">GO20 (Gochain)</option>
                                            <option value="KAVA">KAVA (Kava)</option>
                                            <option value="NEP5">NEP5 (Neo)</option>
                                            <option value="NRC20">NRC20 (Nuls)</option>
                                            <option value="SPL">SPL (Solana)</option>
                                            <option value="TRC21">TRC21 (Tomochain)</option>
                                            <option value="TT20">TT20 (Thundertoken)</option>
                                            <option value="VET">VET (Vechain)</option>
                                            <option value="WAN20">WAN20 (Wanchain)</option>
                                            <option value="POLYGON">POLYGON (Polygon ERC20)</option>
                                        </select>
                                    </td>
                                </tr>
                                <tr>
                                    <td>ID:</td>
                                    <td>
                                        <input v-model="tokenInput.contract" class="input wide"
                                            placeholder="Token contract/ID" size="40"
                                            v-on:change="tokenInputContractChanged()" />
                                        <div class="smallfont padded" v-show="tokenInput.type && !tokenInput.contract">
                                            Search for tokens:
                                            <a v-bind:href="script.assets.explorerUrlForChain(tokenInput.type)">{{script.assets.explorerUrlForChain(tokenInput.type)}}</a>
                                        </div>
                                        <div class="smallfont error padded">
                                            <span v-show="fixedContract"><a v-on:click="tokenInput.contract = fixedContract; tokenInputContractChanged()">Fix</a></span>
                                            {{errorContract}}
                                        </div>
                                    </td>
                                </tr>
                                <tr>
                                    <td>Name:</td>
                                    <td>
                                        <input v-model="tokenInput.name" class="input wide" placeholder="Token name"
                                            v-on:change="tokenInputNameChanged()" />
                                    </td>
                                </tr>
                                <tr>
                                    <td>Symbol:</td>
                                    <td>
                                        <input v-model="tokenInput.symbol" class="input wide" placeholder="Token symbol"
                                            v-on:change="tokenInputSymbolChanged()" />
                                    </td>
                                </tr>
                                <tr>
                                    <td>Decimals:</td>
                                    <td>
                                        <input v-model="tokenInput.decimals" type="number" class="input wide" placeholder="Decimals"
                                            v-on:change="tokenInputDecimalsChanged()" />
                                    </td>
                                </tr>
                                <tr>
                                    <td>Website:</td>
                                    <td>
                                        <input v-model="tokenInput.website" class="input wide"
                                            placeholder="Project website" size="40"
                                            v-on:change="tokenInputWebsiteChanged()" />
                                        <div class="smallfont error padded">
                                            <span v-show="fixedWebsite"><a v-on:click="tokenInput.website = fixedWebsite; tokenInputWebsiteChanged()">Fix</a></span>
                                            {{errorWebsite}}
                                        </div>
                                    </td>
                                </tr>
                                <tr>
                                    <td>Explorer:</td>
                                    <td>
                                        <input v-model="tokenInput.explorerUrl" class="input wide"
                                            placeholder="Token explorer page URL"
                                            size="40" v-on:change="tokenInputExplorerChanged()" />
                                        <div class="smallfont error padded">
                                            <span v-show="fixedExplorer"><a v-on:click="tokenInput.explorerUrl = fixedExplorer; tokenInputExplorerChanged()">Fix</a></span>
                                            {{errorExplorer}}
                                        </div>
                                    </td>
                                </tr>
                                <tr>
                                    <td>Description:</td>
                                    <td>
                                        <textarea v-model="tokenInput.description" class="input wide"
                                            placeholder="Description" rows="2"
                                            v-on:change="tokenInputDescriptionChanged()"></textarea>
                                    </td>
                                </tr>

                                <links-add :tokenInput="tokenInput" />
                                <link-item-edit :tokenInput="tokenInput" v-for="l in this.tokenInput.links" :name="l['name']" :prefix="l['__prefix']"></link-item-edit>
                                <tags-add :tokenInput="tokenInput" />
                            </table>

                            <div class="d-flex justify-content-end">
                                <button class="btn btn-primary m-1" type="button"
                                    v-on:click="checkInputButton()">
                                    {{checkButtonText ? checkButtonText : 'Check'}}
                                </button>
                                <button class="btn btn-outline-primary m-1" type="button" v-on:click="createPullButton()">
                                    {{prButtonText ? prButtonText : 'Create Pull Request'}}
                                </button>
                                <button class="btn btn-outline-primary" type="button" v-show="debugMode"
                                    v-on:click="clearInput()">Clear</button>
                            </div>
                        </div>
                    </div>
                    <div v-show="debugMode" id="add-debug-mode">
                        <div>Debug Mode</div>
                        <div>
                            <button class="btn btn-outline-primary" type="button" v-on:click="debugFillWithDummyData();">Fill
                                test
                                data</button>
                            <input type="checkbox" id="debug_pr_target_fork_cb" v-model="debugPrTargetFork" /><label
                                for="debug_pr_target_fork_cb">Create PR in
                                fork
                            </label>
                        </div>
                    </div>
                </div>
            `
    });

    Vue.component('main-pulls', {
        props: {
            userToken: String,
            loginName: String,
            repo: String,
            enabled: Boolean,
        },
        data: function () {
            return {
                tokenInfo: new script.assets.TokenInfo(),
            }
        },
        methods: {
            onPullSelectToken: async function (token) {
                if (!token) {
                    this.tokenInfo = new script.assets.TokenInfo();
                } else {
                    this.tokenInfo = await script.assets.tokenInfoOfExistingTokenInRepo(token.type, token.id,
                        token.owner, token.repo, token.branch, true, '');
                }
            },
        },
        template:
            `
                <div id="main-pulls" class="flexrow">
                    <div id="pulls-list">
                        <div class="mainheader">Pull Requests</div>
                        <pulls-list :user-token="userToken" :enabled="enabled" v-on:selecttoken="onPullSelectToken"></pulls-list>
                    </div>
                    <token-view :token-info="tokenInfo"></token-view>
                </div>
            `
    });

    Vue.component('main-token-search', {
        props: {
            userToken: String,
            loginName: String,
            repo: String,
        },
        data: function () {
            return {
                tokenInfo: new script.assets.TokenInfo(),
                script: script,
            }
        },
        methods: {
            onSearchSelectToken: async function (token) {
                if (!token || !token.type || !token.asset_id) {
                    this.tokenInfo = new script.assets.TokenInfo();
                } else {
                    const tokenId = script.assets.tokenIdFromAssetId(token.asset_id);
                    this.tokenInfo = await script.assets.tokenInfoOfExistingToken(token.type, tokenId, true, token.asset_id);
                }
            },
        },
        template:
            `
                <div class="row justify-content-center ">
                    <div id="search" class="col-12 col-lg-5 mt-2">
                        <h4>Token Search</h4>
                        <token-search v-on:selecttoken="onSearchSelectToken"></token-search>
                    </div>
                    <token-view :token-info="tokenInfo"></token-view>
                </div>
            `
    });

    Vue.component('warning-no-ghuser', {
        template:
            `
            <div class="row mt-6">
                <div class="col-12">
                    <h3 class="font-weight-bold text-center">Not logged in</h3>
                    <p class="text-center text-gray-400">
                        You need to log in with your GitHub account, and authorize this application.
                    </p>
                </div>
             </div>
            `
    });

    Vue.component('warning-no-forkrepo', {
        props: {
            loginName: String,
            userToken: String
        },
        template:
            `
                <div id="warning-no-forkrepo" class="main center">
                    <div class="mainheader">No fork found</div>
                    <div class="center">
                        <div>No forked repository found for user {{loginName}}</div>
                        <div>You need to fork the assets repository to be able create pull requests.</div>
                        <div class="padded">
                            <button v-on:click="doCreateFork(userToken)">Fork '{{mainRepoName}}' Repo</button>
                        </div>
                        <div>Your fork should be here: <a v-bind:href="gitHub + '/' + loginName + '/' + mainRepoName" target="_blank">{{gitHub}}/{{loginName}}/{{mainRepoName}}</a></div>
                        <div>If you've forked recently, try to <a v-bind:href="window.location">reload</a> after a few minutes.</div>
                        <div>To fork manually, go to the <a v-bind:href="mainRepoUrl" target="_blank">assets repository</a>, and press the Fork button</div>
                        <div>After completion, <a v-bind:href="window.location">reload</a> this page</div>
                        <div class="padded">
                            <img src="img/fork.png" width="600" border="1" />
                        </div>
                    </div>
                </div>
            `
    });

    var app = new Vue({
        el: '#app',
        async created() {
            this.userToken = getUserToken();
            this.loginName = await checkUser(this.userToken);
            this.appMode = getAppMode();
            this.activeTab =
                (this.appMode === 'search') ? 'tab-search' :
                    (this.appMode === 'maintainer') ? 'tab-prs' :
                        'tab-add';
            const resp = await fetch("/get-version");
            if (resp.status == 200) {
                this.version = await resp.text();
            }
            if (this.appMode !== 'search') {
                // check for fork repo; leave it last, slower
                this.repo = await checkRepo(this.userToken, this.loginName, (r) => { this.repoSearchProgress += "."; });
                // Vue v-show-based conditional display causes bad flicker at inital loading, while Vue is compiled;
                // hidden blocks are shown during load.   Use old-fashioned getElementById to show main elements after initialization.
            }
            document.getElementById('main-init').style = "display: block;";
            this.initialized = true;
        },
        data: {
            initialized: false,
            userToken: null,
            loginName: null,
            repo: null,
            version: '',
            activeTab: 'tab-add',
            // appModes:
            // - contributor: add new token, GH login needed
            // - maintainer: all functions, GH login needed
            // - search: search only, no GH login needed
            appMode: 'contributor',
            repoSearchProgress: '.',
            mainRepoName: mainRepoName,
            mainRepoFullName: mainRepoFullName,
            mainRepoUrl: mainRepoUrl,
            gitHub, gitHub,
            window: window,
        },
        methods: {
            selectTab: async function (tab) {
                this.activeTab = tab;
            },
            clearUser: function () {
                this.userToken = null;
                loginName = '';
            },
            logout: function () {
                this.clearUser();
                resetUserToken();
            },
            /*
            toggleMaintainerMode: async function () {
                this.maintainerMode = this.maintainerMode ? false : true;
                window.location.search = updateQueryParams('maintainer', this.maintainerMode ? '1' : '');
            },
            */
            loginActionUrl: function () {
                return assetsAPI;
            },
        }
    });
}
