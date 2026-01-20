export default {
  async fetch(request, env, ctx) {
	 
	//set kv
    await MY_KV_put("myKey","my_val1234")
	
	//get kv
    let value = await MY_KV_get("myKey");
	
    return new Response(`kv test key=myKey,value=${value}`);
  },
};

const root = await navigator.storage.getDirectory();
let KV_space = await root.getDirectoryHandle('tmp', { create: true });

async function MY_KV_put(key,val){
	let MY_KV = await KV_space.getFileHandle(key, { create: true });
	let writable = await MY_KV.createWritable({ keepExistingData: false });
	await writable.write(val);
	await writable.close();
}


async function MY_KV_get(key){
try {
  const MY_KV = await KV_space.getFileHandle(key);
  const data = await MY_KV.getFile();
  return await data.text();
} catch (e) {
  return null
}
}