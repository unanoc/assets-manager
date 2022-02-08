import { TokenInfo, ImageDimensionsCalculator, UrlChecker, LinkItem } from "./tokenInfo";
export declare class TokenInput {
    name: string;
    type: string;
    contract: string;
    symbol: string;
    decimals: number | null;
    logoStream: string;
    logoStreamSize: number;
    logoStreamType: string;
    website: string;
    explorerUrl: string;
    description: string;
    links: LinkItem[];
    tags: string[];
    toTokenInfo(): TokenInfo;
    clone(): TokenInput;
    addLink(name: string, url: string): void;
    addLinkItem(item: LinkItem): void;
}
export declare function checkTokenInput(tokenInput: TokenInput, urlChecker: UrlChecker, imgDimsCalc: ImageDimensionsCalculator, fromBrowser: boolean, checkApiUrl: string): Promise<[number, string, TokenInput | null]>;
export declare function checkTokenInputContract(tokenInput: TokenInput): [number, string, string];
export declare function checkTokenInputWebsite(tokenInput: TokenInput, urlChecker: UrlChecker): Promise<[number, string, string]>;
export declare function checkTokenInputExplorer(tokenInput: TokenInput, urlChecker: UrlChecker): Promise<[number, string, string]>;
export declare function checkTokenInputLogo(tokenInput: TokenInput, imgDimsCalc: ImageDimensionsCalculator): Promise<[number, string]>;
