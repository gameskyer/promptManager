export namespace handlers {
	
	export class AIResponse {
	    success: boolean;
	    data?: any;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new AIResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.data = source["data"];
	        this.error = source["error"];
	    }
	}
	export class AnalyzePromptRequest {
	    prompt: string;
	    config?: services.AIConfig;
	
	    static createFrom(source: any = {}) {
	        return new AnalyzePromptRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.prompt = source["prompt"];
	        this.config = this.convertValues(source["config"], services.AIConfig);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class AtomResponse {
	    success: boolean;
	    data?: any;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new AtomResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.data = source["data"];
	        this.error = source["error"];
	    }
	}
	export class BackupResponse {
	    success: boolean;
	    data?: any;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new BackupResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.data = source["data"];
	        this.error = source["error"];
	    }
	}
	export class BatchImportAtomsRequest {
	    json_data: string;
	
	    static createFrom(source: any = {}) {
	        return new BatchImportAtomsRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.json_data = source["json_data"];
	    }
	}
	export class BuildPromptRequest {
	    atom_ids: number[];
	
	    static createFrom(source: any = {}) {
	        return new BuildPromptRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.atom_ids = source["atom_ids"];
	    }
	}
	export class CategoryResponse {
	    success: boolean;
	    data?: any;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new CategoryResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.data = source["data"];
	        this.error = source["error"];
	    }
	}
	export class CleanupOldVersionsRequest {
	    preset_id: number;
	    keep_count: number;
	
	    static createFrom(source: any = {}) {
	        return new CleanupOldVersionsRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.preset_id = source["preset_id"];
	        this.keep_count = source["keep_count"];
	    }
	}
	export class CreateAtomRequest {
	    value: string;
	    label: string;
	    type: string;
	    category_id: number;
	    synonyms: string[];
	
	    static createFrom(source: any = {}) {
	        return new CreateAtomRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.value = source["value"];
	        this.label = source["label"];
	        this.type = source["type"];
	        this.category_id = source["category_id"];
	        this.synonyms = source["synonyms"];
	    }
	}
	export class CreateCategoryRequest {
	    name: string;
	    type: string;
	    parent_id: number;
	
	    static createFrom(source: any = {}) {
	        return new CreateCategoryRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.type = source["type"];
	        this.parent_id = source["parent_id"];
	    }
	}
	export class CreatePresetRequest {
	    title: string;
	    pos_text: string;
	    neg_text: string;
	    atom_ids: number[];
	    params: Record<string, any>;
	    previews: string[];
	
	    static createFrom(source: any = {}) {
	        return new CreatePresetRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.title = source["title"];
	        this.pos_text = source["pos_text"];
	        this.neg_text = source["neg_text"];
	        this.atom_ids = source["atom_ids"];
	        this.params = source["params"];
	        this.previews = source["previews"];
	    }
	}
	export class CreateVersionRequest {
	    preset_id: number;
	    pos_text: string;
	    neg_text: string;
	    atom_ids: number[];
	    params: Record<string, any>;
	    preview_paths: string[];
	    thumbnail_path: string;
	
	    static createFrom(source: any = {}) {
	        return new CreateVersionRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.preset_id = source["preset_id"];
	        this.pos_text = source["pos_text"];
	        this.neg_text = source["neg_text"];
	        this.atom_ids = source["atom_ids"];
	        this.params = source["params"];
	        this.preview_paths = source["preview_paths"];
	        this.thumbnail_path = source["thumbnail_path"];
	    }
	}
	export class ExplodePromptRequest {
	    prompt: string;
	    config?: services.AIConfig;
	
	    static createFrom(source: any = {}) {
	        return new ExplodePromptRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.prompt = source["prompt"];
	        this.config = this.convertValues(source["config"], services.AIConfig);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ForkPresetRequest {
	    preset_id: number;
	    version_num: number;
	    new_title: string;
	
	    static createFrom(source: any = {}) {
	        return new ForkPresetRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.preset_id = source["preset_id"];
	        this.version_num = source["version_num"];
	        this.new_title = source["new_title"];
	    }
	}
	export class GenericAIRequest {
	    mode: string;
	    prompt: string;
	    config?: services.AIConfig;
	
	    static createFrom(source: any = {}) {
	        return new GenericAIRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.mode = source["mode"];
	        this.prompt = source["prompt"];
	        this.config = this.convertValues(source["config"], services.AIConfig);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class GetAtomsByCategoryRequest {
	    category_id: number;
	    page: number;
	    page_size: number;
	
	    static createFrom(source: any = {}) {
	        return new GetAtomsByCategoryRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.category_id = source["category_id"];
	        this.page = source["page"];
	        this.page_size = source["page_size"];
	    }
	}
	export class GetPresetsRequest {
	    page: number;
	    page_size: number;
	    include_deleted: boolean;
	
	    static createFrom(source: any = {}) {
	        return new GetPresetsRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.page = source["page"];
	        this.page_size = source["page_size"];
	        this.include_deleted = source["include_deleted"];
	    }
	}
	export class ImageResponse {
	    success: boolean;
	    data?: any;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new ImageResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.data = source["data"];
	        this.error = source["error"];
	    }
	}
	export class ImportDataRequest {
	    data: string;
	    merge: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ImportDataRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.data = source["data"];
	        this.merge = source["merge"];
	    }
	}
	export class ImportExtractedRequest {
	    result?: services.ExplodeResult;
	    category_id: number;
	
	    static createFrom(source: any = {}) {
	        return new ImportExtractedRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.result = this.convertValues(source["result"], services.ExplodeResult);
	        this.category_id = source["category_id"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ImportFromJSONRequest {
	    json_data: string;
	
	    static createFrom(source: any = {}) {
	        return new ImportFromJSONRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.json_data = source["json_data"];
	    }
	}
	export class MoveCategoryRequest {
	    id: number;
	    new_parent_id: number;
	
	    static createFrom(source: any = {}) {
	        return new MoveCategoryRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.new_parent_id = source["new_parent_id"];
	    }
	}
	export class OptimizePromptRequest {
	    prompt: string;
	    config?: services.AIConfig;
	
	    static createFrom(source: any = {}) {
	        return new OptimizePromptRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.prompt = source["prompt"];
	        this.config = this.convertValues(source["config"], services.AIConfig);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class PresetResponse {
	    success: boolean;
	    data?: any;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new PresetResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.data = source["data"];
	        this.error = source["error"];
	    }
	}
	export class ReorderCategoriesRequest {
	    ids: number[];
	
	    static createFrom(source: any = {}) {
	        return new ReorderCategoriesRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ids = source["ids"];
	    }
	}
	export class ReverseImagePromptRequest {
	    image_path: string;
	
	    static createFrom(source: any = {}) {
	        return new ReverseImagePromptRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.image_path = source["image_path"];
	    }
	}
	export class SaveAIConfigRequest {
	    config?: services.AIConfig;
	
	    static createFrom(source: any = {}) {
	        return new SaveAIConfigRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.config = this.convertValues(source["config"], services.AIConfig);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class SearchAtomsRequest {
	    search_term: string;
	    type: string;
	    category_id: number;
	    limit: number;
	
	    static createFrom(source: any = {}) {
	        return new SearchAtomsRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.search_term = source["search_term"];
	        this.type = source["type"];
	        this.category_id = source["category_id"];
	        this.limit = source["limit"];
	    }
	}
	export class SearchPresetsRequest {
	    search_term: string;
	    limit: number;
	
	    static createFrom(source: any = {}) {
	        return new SearchPresetsRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.search_term = source["search_term"];
	        this.limit = source["limit"];
	    }
	}
	export class SearchRequest {
	    query: string;
	    limit: number;
	
	    static createFrom(source: any = {}) {
	        return new SearchRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.query = source["query"];
	        this.limit = source["limit"];
	    }
	}
	export class SearchResponse {
	    success: boolean;
	    data?: any;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new SearchResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.data = source["data"];
	        this.error = source["error"];
	    }
	}
	export class SeederResponse {
	    success: boolean;
	    data?: any;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new SeederResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.data = source["data"];
	        this.error = source["error"];
	    }
	}
	export class StarVersionRequest {
	    version_id: number;
	    starred: boolean;
	
	    static createFrom(source: any = {}) {
	        return new StarVersionRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.version_id = source["version_id"];
	        this.starred = source["starred"];
	    }
	}
	export class TranslatePromptRequest {
	    prompt: string;
	    config?: services.AIConfig;
	
	    static createFrom(source: any = {}) {
	        return new TranslatePromptRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.prompt = source["prompt"];
	        this.config = this.convertValues(source["config"], services.AIConfig);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class UpdateAtomRequest {
	    id: number;
	    updates: Record<string, any>;
	
	    static createFrom(source: any = {}) {
	        return new UpdateAtomRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.updates = source["updates"];
	    }
	}
	export class UpdateCategoryRequest {
	    id: number;
	    updates: Record<string, any>;
	
	    static createFrom(source: any = {}) {
	        return new UpdateCategoryRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.updates = source["updates"];
	    }
	}
	export class UpdatePresetRequest {
	    id: number;
	    title: string;
	
	    static createFrom(source: any = {}) {
	        return new UpdatePresetRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.title = source["title"];
	    }
	}
	export class UploadImageRequest {
	    data: string;
	    preset_id: number;
	    version_id: number;
	
	    static createFrom(source: any = {}) {
	        return new UploadImageRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.data = source["data"];
	        this.preset_id = source["preset_id"];
	        this.version_id = source["version_id"];
	    }
	}
	export class VersionResponse {
	    success: boolean;
	    data?: any;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new VersionResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.data = source["data"];
	        this.error = source["error"];
	    }
	}

}

export namespace services {
	
	export class AIConfig {
	    provider: string;
	    api_key: string;
	    endpoint: string;
	    model: string;
	
	    static createFrom(source: any = {}) {
	        return new AIConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.provider = source["provider"];
	        this.api_key = source["api_key"];
	        this.endpoint = source["endpoint"];
	        this.model = source["model"];
	    }
	}
	export class ExtractedAtom {
	    value: string;
	    label: string;
	    type: string;
	    category: string;
	    synonyms: string[];
	    is_new: boolean;
	    existing_id?: number;
	
	    static createFrom(source: any = {}) {
	        return new ExtractedAtom(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.value = source["value"];
	        this.label = source["label"];
	        this.type = source["type"];
	        this.category = source["category"];
	        this.synonyms = source["synonyms"];
	        this.is_new = source["is_new"];
	        this.existing_id = source["existing_id"];
	    }
	}
	export class ExplodeResult {
	    atoms: ExtractedAtom[];
	    raw_prompt: string;
	
	    static createFrom(source: any = {}) {
	        return new ExplodeResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.atoms = this.convertValues(source["atoms"], ExtractedAtom);
	        this.raw_prompt = source["raw_prompt"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

