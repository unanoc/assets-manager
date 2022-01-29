import { isEthereumAddress, toChecksum } from "./eth-address";
import { fetchUniversal, httpPostFromBrowser } from "./fetch";

export const DecimalsMaxValue = 18;
export const InfoAllowedKeys = ["name", "type", "symbol", "decimals", "description", "website", "explorer", "status", "id", "links"];
// Supported keys in links, and their mandatory prefix
export const LinksKeys = {
    // Order matters, appears in this order in app UI
    //"explorer": "",
    "twitter": "https://twitter.com/",
    "github": "https://github.com/",
    "telegram": "https://t.me/",
    "telegram_news": "https://t.me/", // read-only announcement channel
    "blog": "", // blog, other than medium
    "docs": "",
    "forum": "", // community site
    "discord": "https://discord.com/",
    "reddit": "https://reddit.com/",
    "whitepaper": "",
    "medium": "", // url contains 'medium.com'
    "youtube": "https://youtube.com/",
    "facebook": "https://facebook.com/",
    "coinmarketcap": "https://coinmarketcap.com/",
    "coingecko": "https://coingecko.com/",
    "source_code": "" // other than github
};
//const LinksMinRequired = 2;
//const linksKeysString = Object.keys(LinksKeys).reduce(function (agg, item) { return agg + item + ","; }, '');
//const linksMediumContains = 'medium.com';
const assetsAPI = "https://api.assets.trustwallet.com"

export interface TagDescription {
    id: string;
    description: string;
};
export const TagValues: TagDescription[] = [
    {
        id: "stablecoin",
        description: "Tokens that are fixed to an external asset, e.g. the US dollar."
    },
    {
        id: "wrapped",
        description: "Tokens that are wrapped or peg representation of digital assets. Excluded stablecoins"
    },
    {
        id: "synthetics",
        description: "Synthetic assets created to track the value of another asset"
    },
    {
        id: "nft",
        description: "Non-fungible tokens or tokens associated with the NFT ecosystem."
    },
    {
        id: "governance",
        description: "Tokens that used to participate in the governance process of the project."
    },
    {
        id: "defi",
        description: "Tokens that are used for variety of decentralized financial applications."
    },
    {
        id: "staking",
        description: "Tokens that are used for staking to receive rewards."
    },
    {
        id: "staking-native",
        description: "Coins/Blockchains that are used for staking to secure the network to receive rewards."
    },
    {
        id: "privacy",
        description: "Privacy tokens."
    },
    {
        id: "nsfw",
        description: "Content Not suitable for work."
    },
];

export class LinkItem {
    name: string = "";
    url: string = "";
    constructor(name: string, url: string) { this.name = name; this.url = url; }
}

/// Class holding info of a token
export class TokenInfo {
    type: string = "";
    contract: string = "";
    // Logo file path within repo
    logoPath: string = "";
    // Full URL of logo file
    logoUrl: string = "";
    // the logo contents, in base64.  If set, it is used, if empty, Url is used.
    logoStream: string = "";
    logoStreamSize: number = 0;
    logoStreamType: string = "";
    // Info file path within repo
    infoPath: string = "";
    // Full URL of the info file
    infoUrl: string = "";
    info: unknown = {};
    infoString: string = "{}";
    links: LinkItem[] = [];
    tags: string[] = [];
    assetId: string = "";

    explorerUrl(): string {
        return explorerUrlForToken(this.type, this.contract);
    }
}

export function explorerUrlForChain(chainType: string): string {
    switch (chainType.toLowerCase()) {
        case "erc20": return "https://etherscan.io";
        case "trc10":
        case "trc20":
            return "https://tronscan.io";
        case "bep2": return "https://explorer.binance.org";
        case "bep20": return "https://bscscan.com";
        case "thundertoken": return "https://viewblock.io";
    }
    return "";
}

export function explorerUrlForToken(chainType: string, contract: string): string {
    if (contract) {
        switch (chainType.toLowerCase()) {
            case "erc20": return `https://etherscan.io/token/${contract}`;
            case "trc10": return `https://tronscan.io/#/token/${contract}`;
            case "trc20": return `https://tronscan.io/#/token20/${contract}`;
            case "bep2": return `https://explorer.binance.org/asset/${contract}`;
            case "bep20": return `https://bscscan.com/token/${contract}`;
            case "tt20": return `https://viewblock.io/thundercore/address/${contract}`;
            case "nep5": `https://neo.tokenview.com/en/token/0x${contract}`;
            case "nrc20": return `https://nulscan.io/token/info?contractAddress=${contract}`;
            case "vet": return `https://www.wanscan.org/token/${contract}`;
            case "spl": return `https://explorer.solana.com/address/${contract}`;
            case "trc21": return `https://scan.tomochain.com/address/${contract}`;
            case "kava": return "https://www.mintscan.io/kava";
            case "ontology": return "https://explorer.ont.io";
            case "go20": return `https://explorer.gochain.io/addr/${contract}`;
            case "etc20": return `https://blockscout.com/etc/mainnet/tokens/${contract}`;
            case "polygon": return `https://polygonscan.com/token/${contract}`;
        }
    }
    return "";
}

export function chainFromAssetType(tokenType: string) {
    switch (tokenType.toLowerCase()) {
        case "erc20": return "ethereum";
        case "bep2": return "binance";
        case "bep8": return "binance";
        case "bep20": return "smartchain";
        case "trc10": return "tron";
        case "trc20": return "tron";
        case "etc20": return "classic";
        case "thundertoken": return "thundertoken";
        case "wan20": return "wanchain";
        case "trc21": return "tomochain";
        case "tt20": return "thundertoken";
        case "spl": return "solana";
        case "go20": return "gochain";
        case "kava": return "kava";
        case "nep5": return "neo";
        case "nrc20": return "nuls";
        case "vet": return "vechain";
        case "ontology": return "ontology";
        case "polygon": return "polygon";
        default: return "unknown"
    }
}

export function normalizeType(tokenType: string) {
    switch (tokenType.toLowerCase()) {
        case "erc20":
        case "bep2":
        case "bep8":
        case "bep20":
        case "trc10":
        case "trc20":
        case "etc20":
        case "wan20":
        case "trc21":
        case "tt20":
        case "spl":
        case "go20":
        case "kava":
        case "nep5":
        case "nrc20":
        case "vet":
        case "polygon":
            return tokenType.toUpperCase();

        default:
            return ""
    }
}

const mainRepoOwner = "trustwallet";
const mainRepoName = "assets";
const mainMasterBranch = "master";

// Construct TokenInfo for an existing token, specified by type and contract.
export async function tokenInfoOfExistingToken(tokenType: string, contract: string, fetchInfoJson: boolean = false, assetId: string = ''): Promise<TokenInfo> {
    return await tokenInfoOfExistingTokenInRepo(tokenType, contract, mainRepoOwner, mainRepoName, mainMasterBranch, fetchInfoJson, assetId);
}

// Construct TokenInfo for an existing token, specified by type and contract, possible form another repo (typically from a pull request)
export async function tokenInfoOfExistingTokenInRepo(tokenType: string, contract: string, repoOwner: string, repoName: string, branch: string, fetchInfoJson: boolean = false, assetId: string = ''): Promise<TokenInfo> {
    let ti = new TokenInfo();
    ti.type = normalizeType(tokenType);
    ti.contract = contract;
    const chain = chainFromAssetType(tokenType);
    ti.logoPath = `blockchains/${chain}/assets/${ti.contract}/logo.png`;
    ti.logoUrl = `https://raw.githubusercontent.com/${repoOwner}/${repoName}/${branch}/blockchains/${chain}/assets/${ti.contract}/logo.png`;
    ti.infoPath = `blockchains/${chain}/assets/${ti.contract}/info.json`;
    ti.infoUrl = `https://raw.githubusercontent.com/${repoOwner}/${repoName}/${branch}/blockchains/${chain}/assets/${ti.contract}/info.json`;

    ti.logoStream = "";
    ti.logoStreamSize = 0;
    ti.logoStreamType = "";
    ti.info = {};
    ti.infoString = "{}";

    ti.assetId = assetId;

    if (fetchInfoJson) {
        // read info.json
        ti.infoString = "";
        let resp = await fetch(ti.infoUrl);
        if (resp.status == 200) {
            ti.infoString = await resp.text();
            try {
                ti.info = JSON.parse(ti.infoString);
            } catch (error) {
                ti.info = {};
            }
        }
    }

    return ti;
}

// Infer token ID from a logo filename.
// Input: filename, such as "blockchains/ethereum/assets/0x439662426153C4fCB9c6988962FB16475D13d95B/logo.png"
// Output: [type, id], like ["ERC20", "0x439662426153C4fCB9c6988962FB16475D13d95B"]
export function tokenIdFromFile(filename: string): [string, string] {
    const types: string[] = ["ERC20", "BEP2", "BEP20", "TRC10", "TRC20", "ETC20", "WAN20", "TRC21", "TT20", "GO20", "KAVA", "NEP5", "NRC20", "VET", "THUNDERTOKEN", "SPL", "POLYGON"];
    let id: [string, string] = ["", ""];
    types.forEach(type => {
        const chain = chainFromAssetType(type);
        const prefix = `blockchains/${chain}/assets/`;
        const suffix = "/logo.png";
        if (filename.startsWith(prefix)) {
            if (filename.endsWith(suffix)) {
                id = [type, filename.substring(prefix.length, filename.length - suffix.length)];
            }
        }
    });
    // special for TRC10/20 -- both have same 'tron' folder
    if (id[0] && id[0] == "TRC20" && id[1].startsWith("10")) {
        id[0] = "TRC10";
    }
    return id;
}

function isFileAllowed(file: string): boolean {
    // explicit forbidden files
    if (file.endsWith("tokenlist.json")) {
        return false;
    }
    // explicit allowed files
    if (file.startsWith("blockchains") && file.indexOf("assets") > 0) {
        return true;
    }
    if (file.startsWith("blockchains") && file.endsWith("allowlist.json")) {
        return true;
    }
    if (file.startsWith("blockchains") && file.endsWith("validators/list.json")) {
        return true;
    }
    if (file.startsWith("dapps")) {
        return true;
    }
    return false;
}

export class PrFileInfo {
    filename: string = '';
    // added, deleted
    status: string = '';
}

// Check files in the PR
export function checkPrFiles(files: PrFileInfo[]): [number, string] {
    if (files.length === 0) {
        return [2, "No changed files found"];
    }
    const MaxFiles = 20
    if (files.length > MaxFiles) {
        return [2, `Too many changed files, ${files.length}, max ${MaxFiles}.  If not all files are intended, check if forked repo is up to date.`];
    }

    var text = '';
    var result = 0;
    files.forEach(file => {
        if (!isFileAllowed(file.filename)) {
            result = Math.max(result, 2);
            text += `File not allowed: ${file.filename}. Please revert it. \n`;
        }
        if (file.status === 'removed') {
            result = Math.max(result, 2);
            text += `File is being deleted, ${file.filename}. Files should not be deleted in a PR. (Deprecated tokens should be deactivated only.) \n`;
        }
    });
    if (result > 0) {
        return [result, text];
    }
    return [0, `Files OK, ${files.length}`];
}

// Infer token IDs from a logo filenames.
// Input: array with filenames, such as ["blockchains/ethereum/assets/0x439662426153C4fCB9c6988962FB16475D13d95B/logo.png"]
// Output: array of token IDs, [type, id], like [["ERC20", "0x439662426153C4fCB9c6988962FB16475D13d95B"]]
export function tokenIdsFromFiles(filenames: string[]): [string, string][] {
    let ids: [string, string][] = [];
    filenames.forEach(file => {
        const [type, id] = tokenIdFromFile(file);
        if (type && id) {
            ids.push([type, id]);
        }
    });
    return ids;
}

// Plugin for website accesibility check, returns HTTP result code
export interface UrlChecker {
    checkUrl(targetUrl: string): Promise<number>;
}

// Plugin for image dimensions calculator (may use client infra)
export interface ImageDimensionsCalculator {
    get(imageUrl: string, imageStream: string): Promise<{ x: number, y: number }>;
}

export async function checkUrlWithFetch(targetUrl: string): Promise<number> {
    try {
        const result = await fetch(targetUrl);
        if (result.status != 200) {
            return result.status;
        }
        result.status;
    } catch (error) {
        return 404;
    }
}

function errorHead(error): string {
    return error.toString().substring(0, 400 - 1);
}

/*
// return error, warning
function isInfoLinksValid(links: LinkItem[], contract: string, type: string): [string, string] {
    if (!Array.isArray(links)) {
        return [`Links must be an array '${JSON.stringify(links)}' '${contract}' '${type}'`, ""];
    }
    for (let idx = 0; idx < links.length; idx++) {
        const f = links[idx];
        const fname = f['name'];
        if (!fname) {
            return [`Field name missing '${JSON.stringify(f)}'`, ""];
        }
        const furl = f['url'];
        if (!furl) {
            return [`Field url missing '${JSON.stringify(f)}'`, ""];
        }
        const fprefix = f['__prefix'];
        if (fprefix && furl === fprefix) {
            return [`Field url is incomplete '${JSON.stringify(f)}'`, ""];
        }
        // Check there are no other fields
        for (const f2 in f) {
            if (f2 !== 'name' && f2 !== 'url' && f2 !== '__prefix') {
                return [`Invalid field '${f2}' in links '${JSON.stringify(f)}', ${contract} ${type}`, ""];
            }
        }
        if (!Object.prototype.hasOwnProperty.call(LinksKeys, fname)) {
            return [`Not supported field in links '${fname}'.  Supported keys: ${linksKeysString}`, ""];
        }
        const prefix = LinksKeys[fname];
        if (prefix) {
            if (!furl.startsWith(prefix)) {
                return [`Links field '${fname}': '${furl}' must start with '${prefix}'.  Supported keys: ${linksKeysString}`, ""];
            }
        }
        if (!furl.startsWith('https://')) {
            return [`Links field '${fname}': '${furl}' must start with 'https://'.  Supported keys: ${linksKeysString}`, ""];
        }
        // special handling for medium
        if (fname === 'medium') {
            if (!furl.includes(linksMediumContains)) {
                return [`Links field '${fname}': '${furl}' must include '${linksMediumContains}'.  Supported keys: ${linksKeysString}`, ""];
            }
        }
    }

    // check count
    if (!(links && links.length >= LinksMinRequired)) {
        const msg = `At least ${LinksMinRequired} links are required, ${links.length} present.  Add as many as you can: twitter, github, telegram, reddit, etc.`;
        return [msg, ""];
    }    

    return ["", ""];
}
*/

// return error, warning
export function isInfoTagsValid(tags: string[]): [string, string] {
    if (tags.length == 0) {
        return ['At least one tag is needed', ""];
    }
    return ["", ""];
}

// Check tokenInfo for validity: contract is OK, logo is OK, etc.
// returns:
// - result: 0 for all OK, 1 for at least one warning, 2 for at least one error
// - a multi-line string with the detailed results
export async function checkTokenInfo(tokenInfo: TokenInfo, urlChecker: UrlChecker,
    imgDimsCalc: ImageDimensionsCalculator, fromBrowser: boolean): Promise<[number, string]> {
    let res: { res: number, msg: string }[] = [];

    if (!tokenInfo.type || !normalizeType(tokenInfo.type)) {
        res.push({ res: 2, msg: `Invalid token type ${tokenInfo.type}` });
    } else {
        res.push({ res: 0, msg: `Token type OK (${tokenInfo.type})` });
    }

    res.push(checkTokenInfoContract(tokenInfo));

    if (!tokenInfo.info) {
        res.push({ res: 2, msg: "Info.json must not be missing" });
    } else {
        if (fromBrowser) {
            let resp = await httpPostFromBrowser(`${assetsAPI}/v1/validate/asset_info`, tokenInfo.info);
            //console.log(resp);
            if (resp[1]['errors']) {
                for (var k in resp[1]['errors']) {
                    //console.log(resp[1]['errors'][k]);
                    res.push({ res: 2, msg: resp[1]['errors'][k]['message'] });
                }
            }
        }

        // logo
        try {
            (await checkTokenInfoLogo(tokenInfo, imgDimsCalc)).forEach(r => res.push(r));
        } catch (ex) {
            res.push({ res: 1, msg: `Error while checking logo; ${errorHead(ex)}` });
        }
    }

    return AggregateCheckResults(res);
}

// Aggregate results: max and string
export function AggregateCheckResults(res: { res: number, msg: string }[]): [number, string] {
    let maxres = 0;
    let msg = "";
    res.forEach(r => maxres = Math.max(r.res, maxres));
    // Error first, then warnings
    res.forEach(r => { if (r.res >= 2) { msg += "❌  " + r.msg + "\n"; } });
    res.forEach(r => { if (r.res == 1) { msg += "!  " + r.msg + "\n"; } });
    // Finally OKs, all in one line
    let okCount = 0;
    res.forEach(r => {
        if (r.res == 0) {
            if (okCount === 0) { msg += "✓  "; }
            msg += r.msg + "; ";
            okCount++;
        }
    });
    return [maxres, msg];
}

function checkTokenInfoContract(tokenInfo: TokenInfo): { res: number, msg: string } {
    if (!tokenInfo.contract) {
        return { res: 2, msg: "Contract/ID cannot be empty" };
    }
    if (tokenInfo.type.toLowerCase() === "erc20" || tokenInfo.type.toLowerCase() === "bep20") {
        if (!isEthereumAddress(tokenInfo.contract)) {
            return { res: 2, msg: `Contract is not a valid Ethereum address!` };
        }
        const inChecksum = toChecksum(tokenInfo.contract);
        if (inChecksum !== tokenInfo.contract) {
            return { res: 2, msg: `Contract is not in checksum format, should be ${inChecksum} (not ${tokenInfo.contract}). Please rename it. You may need to rename to a temp name first, then to the checksum format, because lowercase-uppercase-only renames are often ignored by the Git client or the filesystem.` };
        }
    }
    return { res: 0, msg: `Contract/ID is OK` };
}

export async function checkTokenInfoLogo(tokenInfo: TokenInfo, imgDimsCalc: ImageDimensionsCalculator): Promise<{ res: number, msg: string }[]> {
    let res: { res: number, msg: string }[] = [];

    if (!tokenInfo.logoStream && !tokenInfo.logoUrl) {
        return [{ res: 2, msg: "Logo image may not be missing" }];
    }
    let logoStreamSize = tokenInfo.logoStreamSize;
    let logoStreamType = tokenInfo.logoStreamType;
    if (!tokenInfo.logoStream) {
        try {
            const response = await fetch(tokenInfo.logoUrl);
            if (response.status != 200) {
                return [{ res: 2, msg: `Could not retrieve logo from url ${tokenInfo.logoUrl}, status ${response.status}` }];
            }
            logoStreamSize = (await response.arrayBuffer()).byteLength;
            logoStreamType = response.headers.get('Content-Type');
        } catch (error) {
            return [{ res: 2, msg: `Could not retrieve logo from url ${tokenInfo.logoUrl}, error ${error}` }];
        }
    }
    if (logoStreamType.toLowerCase() != "image/png") {
        return [{ res: 2, msg: `Logo image must be PNG image (not ${logoStreamType})` }];
    }
    res.push({ res: 0, msg: `Logo image type is OK (${logoStreamType})` });

    if (logoStreamSize > 100000) {
        res.push({ res: 2, msg: `Logo image too large, max 100 kB, current ${logoStreamSize / 1000} kB` });
    } else {
        res.push({ res: 0, msg: `Logo image size is OK (${logoStreamSize / 1000} kB)` });
    }

    try {
        const logoDimension = await imgDimsCalc.get(tokenInfo.logoUrl, tokenInfo.logoStream);
        if (logoDimension.x == 0 && logoDimension.y == 0) {
            res.push({ res: 2, msg: `Could not retrieve logo dimensions` });
        } else if (logoDimension.x > 512 || logoDimension.y > 512) {
            res.push({ res: 2, msg: `Logo should be 256x256 pixels, it is too large ${logoDimension.x}x${logoDimension.y}` });
        } else if (logoDimension.x < 128 || logoDimension.y < 128) {
            res.push({ res: 2, msg: `Logo should be 256x256 pixels, it is too small ${logoDimension.x}x${logoDimension.y}` });
        } else {
            res.push({ res: 0, msg: `Logo dimensions OK (${logoDimension.x}x${logoDimension.y})` });
        }
    } catch (error) {
        res.push({ res: 2, msg: `Could not retrieve logo dimensions (${error})` });
    }

    return res;
}

export async function getExternalTokenInfo(tokenType: string, tokenAddress: string, fromBrowser: boolean): Promise<unknown> {
    try {
        switch (tokenType.toLowerCase()) {
            case 'erc20':
                return await getTokenInfoEthplorer(tokenAddress, fromBrowser);
            case 'bep20':
                return await getTokenInfoBscscan(tokenType, tokenAddress, fromBrowser);
            default:
                // not supported
                //throw `External token info for type ${tokenType} not supported`;
                return null;
        }
    } catch (error) {
        return null;
    }
}

function getTokenCirculationFromExternalInfo(externalTokenInfo: unknown): string {
    return externalTokenInfo["holdersCount"];
}

export async function getTokenCirculation(tokenType: string, tokenAddress: string, fromBrowser: boolean): Promise<string> {
    const externalTokenInfo = await getExternalTokenInfo(tokenType, tokenAddress, fromBrowser);
    //console.log(externalTokenInfo["holdersCount"], externalTokenInfo)
    return getTokenCirculationFromExternalInfo(externalTokenInfo);
}

export async function getTokenCirculationSafe(tokenType: string, tokenAddress: string, fromBrowser: boolean): Promise<string> {
    try {
        const holders = await getTokenCirculation(tokenType, tokenAddress, fromBrowser);
        return holders;
    } catch (ex) {
        console.log('Exception:', ex);
        return '?';
    }
}

async function callEthplorerApi(url: string, fromBrowser: boolean): Promise<unknown> {
    const [status, text] = await fetchUniversal(url, fromBrowser);
    if (status != 200) {
        console.log("ERROR: Non-OK status", status, url);
        return {};
    }
    try {
        return JSON.parse(text);
    }
    catch (error) {
        console.log("error", error)
    }
}

const ethplorerApiUrl = "https://api.ethplorer.io";
const ethplorerApiKey = "freekey";

async function getTokenInfoEthplorer(token: string, fromBrowser: boolean): Promise<unknown> {
    const url = `${ethplorerApiUrl}/getTokenInfo/${token}?apiKey=${ethplorerApiKey}`;
    const data = await callEthplorerApi(url, fromBrowser);
    return {
        symbol: data['symbol'],
        decimals: safeParseInt(data['decimals']),
        holdersCount: data['holdersCount'],
        transfersCount: data['transfersCount'],
        name: data['name'],
        website: data['website'],
        facebook: data['facebook'],
        twitter: data['twitter'],
    }
}

function parseFragment(page: string, fragmentStart: string, fragmentEnd: string, url: string): string {
    const idx1 = page.indexOf(fragmentStart);
    const rangeLen = 20;
    if (idx1 < 1) {
        throw `Could not parse item from explorer page; opening fragment not found, ${url} '${fragmentStart}' ${idx1} ${page.length} ${errorHead(page)}`;
    }
    const range = page.substring(idx1 + fragmentStart.length, idx1 + fragmentStart.length + rangeLen);
    //console.log('range', range); 
    const idx2 = range.indexOf(fragmentEnd);
    //console.log("idx2", idx2);
    if (idx2 < 0) {
        throw `Could not parse item from explorer page; closing fragment not found, ${url} '${fragmentEnd}' ${idx1} ${idx2} ${page.length} ${range}`;
    }
    let fragment = page.substring(idx1 + fragmentStart.length, idx1 + fragmentStart.length + idx2);
    //console.log(`Fragment between '${fragmentStart}' and '${fragmentEnd}': '${fragment}' (${page.length} ${idx1} ${idx2} ${range} ${rangeLen})`);
    return fragment;
}

function parseFragmentFromEnd(page: string, fragmentStart: string, fragmentEnd: string, url: string): string {
    const idx2 = page.indexOf(fragmentEnd);
    const rangeLen = 20;
    if (idx2 < rangeLen) {
        throw `Could not parse item from explorer page; closing fragment not found, ${url} '${fragmentEnd}' ${idx2} ${page.length} ${errorHead(page)}`;
    }
    const range = page.substring(idx2 - rangeLen, idx2);
    const idx1 = range.indexOf(fragmentStart);
    //console.log("idx1", idx1);
    if (idx1 < 0) {
        throw `Could not parse item from explorer page; opening fragment not found, ${url} '${fragmentStart}' ${idx1} ${idx2} ${page.length} ${range}`;
    }
    let fragment = page.substring(idx2 - rangeLen + idx1 + fragmentStart.length, idx2);
    //console.log(`Fragment between '${fragmentStart}' and '${fragmentEnd}': '${fragment}' (${page.length} ${idx1} ${idx2} ${range} ${rangeLen})`);
    return fragment;
}

async function getTokenInfoBscscan(tokenType: string, tokenAddress: string, fromBrowser: boolean) {
    const explorerUrl = explorerUrlForToken(tokenType, tokenAddress);
    //console.log("explorerUrl", explorerUrl);
    const [status, text] = await fetchUniversal(explorerUrl, fromBrowser);
    if (status != 200) {
        throw `Could not retrieve explorer page ${explorerUrl} ${status} ${errorHead(text)}`;
    }

    var symbol = '';
    var decimalsString = '';
    var holdersString = '';
    try {
        const textPre = text.replace('\n', '').replace('\r', '').replace('\\n', '').replace('\\r', '');

        try {
            symbol = parseFragment(textPre, "'symbol': $.sanitize('", "'", explorerUrl);
        } catch (ex) {
            console.log(`Exception, ignored; ${ex}`);
        }
        if (!symbol) {
            try {
                symbol = parseFragment(textPre, "'symbol': '", "'", explorerUrl);
            } catch (ex) {
                console.log(`Exception, ignored; ${ex}`);
            }
        }

        try {
            decimalsString = parseFragment(textPre, "decimals': '", "'", explorerUrl).replace('\n', '').replace('\\n', '');
        } catch (ex) {
            console.log(`Exception, ignored; ${ex}`);
        }

        try {
            holdersString = parseFragmentFromEnd(textPre, '">', ' addresses', explorerUrl).replace(',', '').replace('.', '').replace('\n', '').replace('\\n', '');
        } catch (ex) {
            console.log(`Exception, ignored; ${ex}`);
        }
    } catch (ex) {
        console.log(`Exception, ignored; ${ex}`);
    }

    //console.log(explorerUrl, holdersString);
    const externalTokenInfo = {
        symbol: symbol,
        decimals: decimalsString,
        holdersCount: holdersString,
    }
    //console.log('externalTokenInfo', externalTokenInfo);
    return externalTokenInfo;
}

const CirculationHoldersLimit = 25000;

/*
async function retrieveAndCheckHoldersLimit(tokenType: string, tokenAddress: string, fromBrowser: boolean): Promise<{res: number, msg: string}> {
    try {
        const externalTokenInfo = await getExternalTokenInfo(tokenType, tokenAddress, fromBrowser);
        return checkHoldersLimit(externalTokenInfo);
    } catch (ex) {
        return {res: 1, msg: `No. of holders could not be checked; ${errorHead(ex)}`};
    }
}
*/

function checkHoldersLimit(externalTokenInfo: unknown): { res: number, msg: string } {
    if (!externalTokenInfo) {
        return { res: 1, msg: `No. of holders not checked, external info not available; ${externalTokenInfo}` };
    }
    const holders = getTokenCirculationFromExternalInfo(externalTokenInfo);
    if (!holders || holders === '?') {
        // could not check
        return { res: 1, msg: `No. of holders not checked; '${holders}'` };
    }
    var holdersNum = safeParseInt(holders);
    //console.log('holders', holders, holdrsNum);
    if (holdersNum == NaN) {
        return { res: 1, msg: `No. of holders not checked, NaN; (${holdersNum} ${holders})` };
    }

    if (holdersNum >= CirculationHoldersLimit) {
        return { res: 0, msg: `Token circulation OK (no. of holders: ${holdersNum})` };
    }
    return { res: 2, msg: `Low token circulation: no. of holders is ${holdersNum}, below limit of ${CirculationHoldersLimit}` };
}

export function safeParseInt(value: string): number {
    try {
        const num: number = parseInt(value, 10);
        if (num === NaN || !num) {
            return 0;
        }
        return num;
    } catch (err) {
        return 0;
    }
}

/// Parse out Token ID from Asset ID.  E.g. "c60_t0x8eEF5a82E6Aa222a60F009ac18c24EE12dBf4b41"
export function tokenIdFromAssetId(assetId: string): string {
    const separator = '_t';
    const idx = assetId.indexOf(separator);
    if (idx < 0) {
        return assetId;
    }
    return assetId.substring(idx + separator.length, assetId.length);
}