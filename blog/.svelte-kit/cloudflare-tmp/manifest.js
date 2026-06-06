export const manifest = (() => {
function __memo(fn) {
	let value;
	return () => value ??= (value = fn());
}

return {
	appDir: "_app",
	appPath: "_app",
	assets: new Set([".gitkeep"]),
	mimeTypes: {},
	_: {
		client: {start:"_app/immutable/entry/start.aETX-UfI.js",app:"_app/immutable/entry/app.CA7gUK1n.js",imports:["_app/immutable/entry/start.aETX-UfI.js","_app/immutable/chunks/bIT_w4Yp.js","_app/immutable/chunks/Cvc1K2yR.js","_app/immutable/entry/app.CA7gUK1n.js","_app/immutable/chunks/Cvc1K2yR.js","_app/immutable/chunks/kNaey6uv.js","_app/immutable/chunks/xihTtKlq.js"],stylesheets:[],fonts:[],uses_env_dynamic_public:false},
		nodes: [
			__memo(() => import('../output/server/nodes/0.js')),
			__memo(() => import('../output/server/nodes/1.js')),
			__memo(() => import('../output/server/nodes/4.js'))
		],
		remotes: {
			
		},
		routes: [
			{
				id: "/about",
				pattern: /^\/about\/?$/,
				params: [],
				page: { layouts: [0,], errors: [1,], leaf: 2 },
				endpoint: null
			},
			{
				id: "/api/visit",
				pattern: /^\/api\/visit\/?$/,
				params: [],
				page: null,
				endpoint: __memo(() => import('../output/server/entries/endpoints/api/visit/_server.ts.js'))
			}
		],
		prerendered_routes: new Set(["/","/__data.json","/conways-law-at-level-5","/conways-law-at-level-5/__data.json","/wardley-was-right","/wardley-was-right/__data.json","/i-followed-the-problem-home","/i-followed-the-problem-home/__data.json","/treating-ideas-as-releasable-software","/treating-ideas-as-releasable-software/__data.json"]),
		matchers: async () => {
			
			return {  };
		},
		server_assets: {}
	}
}
})();

export const prerendered = new Set(["/","/__data.json","/conways-law-at-level-5","/conways-law-at-level-5/__data.json","/wardley-was-right","/wardley-was-right/__data.json","/i-followed-the-problem-home","/i-followed-the-problem-home/__data.json","/treating-ideas-as-releasable-software","/treating-ideas-as-releasable-software/__data.json"]);

export const base_path = "";
