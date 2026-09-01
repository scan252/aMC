import { createServer } from "node:http";
import { readFile } from "node:fs/promises";
import { join, extname } from "node:path";
const root = "D:/Zcode/git/aMC/design";
const mime = {".html":"text/html; charset=utf-8",".css":"text/css",".js":"text/javascript",".png":"image/png",".svg":"image/svg+xml"};
createServer(async (req,res)=>{
  try{
    let p = req.url.split("?")[0]; if(p==="/") p="/index.html";
    const data = await readFile(join(root, decodeURIComponent(p)));
    res.writeHead(200,{"Content-Type":mime[extname(p)]||"application/octet-stream"});
    res.end(data);
  }catch(e){ res.writeHead(404); res.end("nf"); }
}).listen(8123,"127.0.0.1",()=>console.log("serving on 8123"));
