declare global {
	namespace App {
		interface Platform {
			env: {
				BLOG_ANALYTICS: KVNamespace;
			};
		}
	}
}

export {};
