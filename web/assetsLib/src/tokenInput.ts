import {
    TokenInfo,
    ImageDimensionsCalculator,
    explorerUrlForToken,
    AggregateCheckResults,
    UrlChecker,
    LinkItem,
    checkTokenInfo,
    checkUrlWithFetch,
    checkTokenInfoLogo,
} from "./tokenInfo";
import { isEthereumAddress, toChecksum } from "./eth-address";

/// Class for entering input for a token
export class TokenInput {
    name: string = "";
    type: string = "ERC20";
    contract: string = "";
    symbol: string = "";
    decimals: number|null = null;
    // base64-encoded logo image stream
    logoStream: string = "";
    logoStreamSize: number = 0;
    logoStreamType: string = "";
    website: string = "";
    explorerUrl: string = "";
    description: string = "";
    links: LinkItem[] = [];
    tags: string[] = [];

    toTokenInfo() {
        let tokenInfo = new TokenInfo();
        tokenInfo.type = this.type;
        tokenInfo.contract = this.contract;
        tokenInfo.logoUrl = "";
        tokenInfo.logoStream = this.logoStream;
        tokenInfo.logoStreamSize = this.logoStreamSize;
        tokenInfo.logoStreamType = this.logoStreamType;
        tokenInfo.infoUrl = "";
        tokenInfo.info = {
            name: this.name,
            type: this.type.toUpperCase(),
            symbol: this.symbol,
            decimals: this.decimals || 0,
            website: this.website,
            description: this.description,
            explorer: this.explorerUrl,
            status: "active",
            id: this.contract,
            links: [],
            tags: this.tags,
        };
        // copy links (to leave out addl internal fields)
        for (let i in this.links) {
            tokenInfo.info['links'].push(new LinkItem(this.links[i].name, this.links[i].url));
        }
        tokenInfo.infoString = JSON.stringify(tokenInfo.info, null, 4);
        tokenInfo.links = this.links;
        tokenInfo.tags = this.tags;
        return tokenInfo;
    }

    clone(): TokenInput {
        let ti2 = new TokenInput();
        ti2.name = this.name;
        ti2.type = this.type;
        ti2.contract = this.contract;
        ti2.logoStream = this.logoStream;
        ti2.logoStreamSize = this.logoStreamSize;
        ti2.logoStreamType = this.logoStreamType;
        ti2.website = this.website;
        ti2.explorerUrl = this.explorerUrl;
        ti2.description = this.description;
        ti2.links = this.links;
        return ti2;
    }

    addLink(name: string, url: string): void {
        this.addLinkItem(new LinkItem(name, url));
    }
    addLinkItem(item: LinkItem): void {
        if (!item.name) { return; }
        if (this.links.find(x => (x.name === item.name)) !== undefined) { return; }  // already present
        this.links.push(item);
    }
}

// Check tokenInput for validity: everything is filled, logo is OK, etc.
// returns:
// - result: 0 for all OK, 1 for at least on warning, 2 for at least on error
// - a multi-line string with the detailed results
// - fixed version if can be auto-fixed
export async function checkTokenInput(tokenInput: TokenInput, urlChecker: UrlChecker, imgDimsCalc: ImageDimensionsCalculator, fromBrowser: boolean, checkApiUrl: string): Promise<[number, string, TokenInput | null]> {
    let res: { res: number, msg: string }[] = [];
    let fixed: TokenInput | null = null;

    if (!tokenInput.name) {
        res.push({ res: 2, msg: "Name cannot be empty" });
    } else {
        res.push({ res: 0, msg: "Name is set" });
    }

    // logo must be loaded (stream)
    if (!tokenInput.logoStream || tokenInput.logoStream.length < 10) {
        res.push({ res: 2, msg: "Logo image may not be missing"});
    }

    // convert to tokenInfo, check that
    const tokenInfo = tokenInput.toTokenInfo();
    const [resnumTI, resmsgTI] = await checkTokenInfo(tokenInfo, { checkUrl: checkUrlWithFetch }, imgDimsCalc, true, checkApiUrl);
    res.push({ res: resnumTI, msg: resmsgTI });

    const [resnum, resmsg] = AggregateCheckResults(res);
    return [resnum, resmsg, fixed];
}

export function checkTokenInputContract(tokenInput: TokenInput): [number, string, string] {
    if (!tokenInput.contract) {
        return [2, "Contract/ID cannot be empty", null];
    }
    if (tokenInput.type.toLowerCase() === "erc20" || tokenInput.type.toLowerCase() === "bep20") {
        if (!isEthereumAddress(tokenInput.contract)) {
            return [2, `Contract is not a valid Ethereum address!`, null];
        }
        const inChecksum = toChecksum(tokenInput.contract);
        if (inChecksum !== tokenInput.contract) {
            return [2, `Contract is not in checksum format, should be ${inChecksum}`, inChecksum];
        }
    }
    return [0, `Contract/ID is OK`, null];
}

export async function checkTokenInputWebsite(tokenInput: TokenInput, urlChecker: UrlChecker): Promise<[number, string, string]> {
    if (!tokenInput.website) {
        return [2, "Website cannot be empty", null];
    }
    var website = tokenInput.website;
    // should start with http
    const prefix = 'https://';
    if (!website.startsWith(prefix.substring(0, website.length))) {
        const fixed = prefix + website;
        return [2, `Website should start with '${prefix}', ${website}`, fixed];
    }
    try {
        const result = await urlChecker.checkUrl(website);
        if (result >= 400 && result <= 499) {
            return [2, `Website does not exist, status ${result}, url ${website}`, null];
        }
        if (result != 200) {
            return [1, `Could not check website availability, status ${result}, url ${website}`, null];
        }
    } catch (error) {
        // may be CORS, treat only as warning
        return [1, `Could not check website availability, error ${error}, url ${website}`, null];
    }
    return [0, `Website OK`, null];
}

export async function checkTokenInputExplorer(tokenInput: TokenInput, urlChecker: UrlChecker): Promise<[number, string, string]> {
    if (!tokenInput.explorerUrl) {
        return [2, "Explorer cannot be empty", null];
    }
    const explorer = tokenInput.explorerUrl;
    try {
        const result = await urlChecker.checkUrl(explorer);
        if (result == 404) {
            return [2, `ExplorerUrl does not exist, status ${result}, url ${explorer}`, null];
        }
        if (result != 200) {
            return [1, `Could not check if ExplorerUrl exists, status ${result}, url ${tokenInput.explorerUrl}`, null];
        }
        // check if explorer is what we would think
        const guessedExplorer = explorerUrlForToken(tokenInput.type, tokenInput.contract);
        if (explorer != guessedExplorer) {
            return [1, `Recommended ExplorerUrl is ${guessedExplorer} instead of ${explorer}`, guessedExplorer];
        }
    } catch (error) {
        // may be CORS, treat only as warning
        return [1, `Could not check if ExplorerUrl exists, error ${error}, url ${explorer}`, null];
    }
    return [0, `ExplorerUrl OK`, null];
}

async function checkTokenInputLogoInternal(tokenInput: TokenInput, imgDimsCalc: ImageDimensionsCalculator): Promise<[number, string][]> {

    if (!tokenInput.logoStream || tokenInput.logoStream.length < 10) {
        return [[2, "Logo image may not be missing"]];
    }

    // convert to tokenInfo, check that
    var tokenInfo = new TokenInfo();
    tokenInfo.logoStream = tokenInput.logoStream;
    tokenInfo.logoStreamSize = tokenInput.logoStreamSize;
    tokenInfo.logoStreamType = tokenInput.logoStreamType;

    const res1: { res: number, msg: string }[] = await checkTokenInfoLogo(tokenInfo, imgDimsCalc);
    var res2: [number, string][] = [];
    for (var r in res1) {
        res2.push([res1[r].res, res1[r].msg]);
    }
    return res2;
}

export async function checkTokenInputLogo(tokenInput: TokenInput, imgDimsCalc: ImageDimensionsCalculator): Promise<[number, string]> {
    const logoRes = await checkTokenInputLogoInternal(tokenInput, imgDimsCalc);
    let res2: { res: number, msg: string }[] = [];
    logoRes.forEach((r) => res2.push({res: r[0], msg: r[1]}));
    const [resnum, resmsg] = AggregateCheckResults(res2);
    return [resnum, resmsg];
}
