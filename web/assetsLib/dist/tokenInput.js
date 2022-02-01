"use strict";
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.checkTokenInputLogo = exports.checkTokenInputExplorer = exports.checkTokenInputWebsite = exports.checkTokenInputContract = exports.checkTokenInput = exports.TokenInput = void 0;
const tokenInfo_1 = require("./tokenInfo");
const eth_address_1 = require("./eth-address");
/// Class for entering input for a token
class TokenInput {
    constructor() {
        this.name = "";
        this.type = "ERC20";
        this.contract = "";
        this.symbol = "";
        this.decimals = null;
        // base64-encoded logo image stream
        this.logoStream = "";
        this.logoStreamSize = 0;
        this.logoStreamType = "";
        this.website = "";
        this.explorerUrl = "";
        this.description = "";
        this.links = [];
        this.tags = [];
    }
    toTokenInfo() {
        let tokenInfo = new tokenInfo_1.TokenInfo();
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
            tokenInfo.info['links'].push(new tokenInfo_1.LinkItem(this.links[i].name, this.links[i].url));
        }
        tokenInfo.infoString = JSON.stringify(tokenInfo.info, null, 4);
        tokenInfo.links = this.links;
        tokenInfo.tags = this.tags;
        return tokenInfo;
    }
    clone() {
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
    addLink(name, url) {
        this.addLinkItem(new tokenInfo_1.LinkItem(name, url));
    }
    addLinkItem(item) {
        if (!item.name) {
            return;
        }
        if (this.links.find(x => (x.name === item.name)) !== undefined) {
            return;
        } // already present
        this.links.push(item);
    }
}
exports.TokenInput = TokenInput;
// Check tokenInput for validity: everything is filled, logo is OK, etc.
// returns:
// - result: 0 for all OK, 1 for at least on warning, 2 for at least on error
// - a multi-line string with the detailed results
// - fixed version if can be auto-fixed
function checkTokenInput(tokenInput, urlChecker, imgDimsCalc, fromBrowser) {
    return __awaiter(this, void 0, void 0, function* () {
        let res = [];
        let fixed = null;
        if (!tokenInput.name) {
            res.push({ res: 2, msg: "Name cannot be empty" });
        }
        else {
            res.push({ res: 0, msg: "Name is set" });
        }
        // logo must be loaded (stream)
        if (!tokenInput.logoStream || tokenInput.logoStream.length < 10) {
            res.push({ res: 2, msg: "Logo image may not be missing" });
        }
        // convert to tokenInfo, check that
        const tokenInfo = tokenInput.toTokenInfo();
        const [resnumTI, resmsgTI] = yield tokenInfo_1.checkTokenInfo(tokenInfo, { checkUrl: tokenInfo_1.checkUrlWithFetch }, imgDimsCalc, true);
        res.push({ res: resnumTI, msg: resmsgTI });
        const [resnum, resmsg] = tokenInfo_1.AggregateCheckResults(res);
        return [resnum, resmsg, fixed];
    });
}
exports.checkTokenInput = checkTokenInput;
function checkTokenInputContract(tokenInput) {
    if (!tokenInput.contract) {
        return [2, "Contract/ID cannot be empty", null];
    }
    if (tokenInput.type.toLowerCase() === "erc20" || tokenInput.type.toLowerCase() === "bep20") {
        if (!eth_address_1.isEthereumAddress(tokenInput.contract)) {
            return [2, `Contract is not a valid Ethereum address!`, null];
        }
        const inChecksum = eth_address_1.toChecksum(tokenInput.contract);
        if (inChecksum !== tokenInput.contract) {
            return [2, `Contract is not in checksum format, should be ${inChecksum}`, inChecksum];
        }
    }
    return [0, `Contract/ID is OK`, null];
}
exports.checkTokenInputContract = checkTokenInputContract;
function checkTokenInputWebsite(tokenInput, urlChecker) {
    return __awaiter(this, void 0, void 0, function* () {
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
            const result = yield urlChecker.checkUrl(website);
            if (result >= 400 && result <= 499) {
                return [2, `Website does not exist, status ${result}, url ${website}`, null];
            }
            if (result != 200) {
                return [1, `Could not check website availability, status ${result}, url ${website}`, null];
            }
        }
        catch (error) {
            // may be CORS, treat only as warning
            return [1, `Could not check website availability, error ${error}, url ${website}`, null];
        }
        return [0, `Website OK`, null];
    });
}
exports.checkTokenInputWebsite = checkTokenInputWebsite;
function checkTokenInputExplorer(tokenInput, urlChecker) {
    return __awaiter(this, void 0, void 0, function* () {
        if (!tokenInput.explorerUrl) {
            return [2, "Explorer cannot be empty", null];
        }
        const explorer = tokenInput.explorerUrl;
        try {
            const result = yield urlChecker.checkUrl(explorer);
            if (result == 404) {
                return [2, `ExplorerUrl does not exist, status ${result}, url ${explorer}`, null];
            }
            if (result != 200) {
                return [1, `Could not check if ExplorerUrl exists, status ${result}, url ${tokenInput.explorerUrl}`, null];
            }
            // check if explorer is what we would think
            const guessedExplorer = tokenInfo_1.explorerUrlForToken(tokenInput.type, tokenInput.contract);
            if (explorer != guessedExplorer) {
                return [1, `Recommended ExplorerUrl is ${guessedExplorer} instead of ${explorer}`, guessedExplorer];
            }
        }
        catch (error) {
            // may be CORS, treat only as warning
            return [1, `Could not check if ExplorerUrl exists, error ${error}, url ${explorer}`, null];
        }
        return [0, `ExplorerUrl OK`, null];
    });
}
exports.checkTokenInputExplorer = checkTokenInputExplorer;
function checkTokenInputLogoInternal(tokenInput, imgDimsCalc) {
    return __awaiter(this, void 0, void 0, function* () {
        if (!tokenInput.logoStream || tokenInput.logoStream.length < 10) {
            return [[2, "Logo image may not be missing"]];
        }
        // convert to tokenInfo, check that
        var tokenInfo = new tokenInfo_1.TokenInfo();
        tokenInfo.logoStream = tokenInput.logoStream;
        tokenInfo.logoStreamSize = tokenInput.logoStreamSize;
        tokenInfo.logoStreamType = tokenInput.logoStreamType;
        const res1 = yield tokenInfo_1.checkTokenInfoLogo(tokenInfo, imgDimsCalc);
        var res2 = [];
        for (var r in res1) {
            res2.push([res1[r].res, res1[r].msg]);
        }
        return res2;
    });
}
function checkTokenInputLogo(tokenInput, imgDimsCalc) {
    return __awaiter(this, void 0, void 0, function* () {
        const logoRes = yield checkTokenInputLogoInternal(tokenInput, imgDimsCalc);
        let res2 = [];
        logoRes.forEach((r) => res2.push({ res: r[0], msg: r[1] }));
        const [resnum, resmsg] = tokenInfo_1.AggregateCheckResults(res2);
        return [resnum, resmsg];
    });
}
exports.checkTokenInputLogo = checkTokenInputLogo;
