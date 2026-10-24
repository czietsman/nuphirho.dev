// cloudflare:workers exists only in the Workers runtime. Tests resolve it to
// this empty environment, and a test that needs a binding mocks it.
export const env = {};
