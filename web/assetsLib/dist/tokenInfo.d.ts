export declare const DecimalsMaxValue = 18;
export declare const InfoAllowedKeys: string[];
export declare const LinksKeys: {
    twitter: string;
    github: string;
    telegram: string;
    telegram_news: string;
    blog: string;
    docs: string;
    forum: string;
    discord: string;
    reddit: string;
    whitepaper: string;
    medium: string;
    youtube: string;
    facebook: string;
    coinmarketcap: string;
    coingecko: string;
    source_code: string;
};
export interface TagDescription {
    id: string;
    description: string;
}
export declare const TagValues: TagDescription[];
export declare class LinkItem {
    name: string;
    url: string;
    constructor(name: string, url: string);
}
export declare class TokenInfo {
    type: string;
    contract: string;
    logoPath: string;
    logoUrl: string;
    logoStream: string;
    logoStreamSize: number;
    logoStreamType: string;
    infoPath: string;
    infoUrl: string;
    info: unknown;
    infoString: string;
    links: LinkItem[];
    tags: string[];
    assetId: string;
    explorerUrl(): string;
}
export declare function explorerUrlForChain(chainType: string): string;
export declare function explorerUrlForToken(chainType: string, contract: string): string;
export declare function chainFromAssetType(tokenType: string): "classic" | "polygon" | "unknown" | "ethereum" | "thundertoken" | "kava" | "ontology" | "binance" | "smartchain" | "tron" | "wanchain" | "tomochain" | "solana" | "gochain" | "neo" | "nuls" | "vechain";
export declare function normalizeType(tokenType: string): string;
export declare function tokenInfoOfExistingToken(tokenType: string, contract: string, fetchInfoJson?: boolean, assetId?: string): Promise<TokenInfo>;
export declare function tokenInfoOfExistingTokenInRepo(tokenType: string, contract: string, repoOwner: string, repoName: string, branch: string, fetchInfoJson?: boolean, assetId?: string): Promise<TokenInfo>;
export declare function tokenIdFromFile(filename: string): [string, string];
export declare class PrFileInfo {
    filename: string;
    status: string;
}
export declare function checkPrFiles(files: PrFileInfo[]): [number, string];
export declare function tokenIdsFromFiles(filenames: string[]): [string, string][];
export interface UrlChecker {
    checkUrl(targetUrl: string): Promise<number>;
}
export interface ImageDimensionsCalculator {
    get(imageUrl: string, imageStream: string): Promise<{
        x: number;
        y: number;
    }>;
}
export declare function checkUrlWithFetch(targetUrl: string): Promise<number>;
export declare function isInfoTagsValid(tags: string[]): [string, string];
export declare function checkTokenInfo(tokenInfo: TokenInfo, urlChecker: UrlChecker, imgDimsCalc: ImageDimensionsCalculator, fromBrowser: boolean): Promise<[number, string]>;
export declare function AggregateCheckResults(res: {
    res: number;
    msg: string;
}[]): [number, string];
export declare function checkTokenInfoLogo(tokenInfo: TokenInfo, imgDimsCalc: ImageDimensionsCalculator): Promise<{
    res: number;
    msg: string;
}[]>;
export declare function getExternalTokenInfo(tokenType: string, tokenAddress: string, fromBrowser: boolean): Promise<unknown>;
export declare function getTokenCirculation(tokenType: string, tokenAddress: string, fromBrowser: boolean): Promise<string>;
export declare function getTokenCirculationSafe(tokenType: string, tokenAddress: string, fromBrowser: boolean): Promise<string>;
export declare function safeParseInt(value: string): number;
export declare function tokenIdFromAssetId(assetId: string): string;
