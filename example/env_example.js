export default {
  async fetch(request, env, ctx) {
	const value = env.API_URL
    return new Response(`API_URL = ${value}`)
  },
};