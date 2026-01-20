package main
var datea=`
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Web IDE Interface</title>
    <style> *{box-sizing:border-box;margin:0;padding:0}body{font-family:'Segoe UI',Tahoma,Geneva,Verdana,sans-serif;height:100vh;overflow:hidden;background-color:#1e1e1e;color:#d4d4d4;display:flex}.editor-container{flex:3;border-right:1px solid #333;display:flex;flex-direction:column}.editor-header{background-color:#252526;padding:10px 20px;font-size:14px;color:#9cdcfe;border-bottom:1px solid #333}#code-input{flex:1;width:100%;background-color:#1e1e1e;color:#d4d4d4;border:none;outline:none;padding:20px;font-family:'Consolas','Courier New',monospace;font-size:14px;line-height:1.5;resize:none;white-space:pre}.sidebar{flex:1;background-color:#252526;display:flex;flex-direction:column;padding:20px;gap:20px}h3{font-size:16px;color:#fff;margin-bottom:10px;border-bottom:1px solid #3e3e42;padding-bottom:5px}#run-btn{background-color:#0e639c;color:white;border:none;padding:12px;border-radius:4px;cursor:pointer;font-size:16px;font-weight:bold;transition:background-color 0.2s;width:100%}#run-btn:hover{background-color:#1177bb}#run-btn:active{background-color:#094770}.env-section{flex:1;display:flex;flex-direction:column;overflow:hidden}.env-list{flex:1;overflow-y:auto;margin-bottom:10px;padding-right:5px}.env-list::-webkit-scrollbar{width:8px}.env-list::-webkit-scrollbar-thumb{background:#424242;border-radius:4px}.env-row{display:flex;gap:5px;margin-bottom:8px}.env-input{background-color:#3c3c3c;border:1px solid #3c3c3c;color:#cccccc;padding:6px;border-radius:3px;width:100%;font-size:12px}.env-input:focus{border-color:#0e639c;outline:none}.icon-btn{background:none;border:1px solid #454545;color:#cccccc;cursor:pointer;width:30px;border-radius:3px;display:flex;align-items:center;justify-content:center}.icon-btn:hover{background-color:#3e3e42}.delete-btn:hover{background-color:#a12626;border-color:#a12626}#add-env-btn{background-color:#3e3e42;color:white;border:none;padding:8px;border-radius:4px;cursor:pointer;width:100%;font-size:12px}#add-env-btn:hover{background-color:#505055}.modal-overlay{position:fixed;inset:0;background:rgba(0,0,0,0.6);display:none;align-items:center;justify-content:center;z-index:999}.modal{background-color:#252526;border:1px solid #3e3e42;border-radius:6px;padding:20px;width:320px;box-shadow:0 10px 30px rgba(0,0,0,0.5)}.modal h4{margin-bottom:10px;color:#9cdcfe;font-size:16px}.modal p{font-size:14px;margin-bottom:20px}.modal-actions{display:flex;gap:10px;justify-content:flex-end}.modal-btn{padding:8px 14px;border-radius:4px;border:none;cursor:pointer;font-size:13px}.btn-close{background-color:#3e3e42;color:#fff}.btn-view{background-color:#0e639c;color:#fff}
    </style>
</head>
<body>
    <div class="editor-container">
        <div class="editor-header">main.js</div>
        <textarea id="code-input" spellcheck="false" placeholder="// 在此处输入您的代码...">export default {
  async fetch(request, env, ctx) {
    return new Response("Hello world!");
  },
};</textarea>
    </div>
    <div class="sidebar">
        <div>
            <h3>操作</h3>
            <button id="run-btn">▷ 立即运行</button>
        </div>
        <div class="env-section">
            <h3>设置环境变量</h3>
            <div id="env-list" class="env-list">
                <div class="env-row">
                    <input type="text" class="env-input key" placeholder="KEY">
                    <input type="text" class="env-input val" placeholder="VALUE">
                    <button class="icon-btn delete-btn" onclick="removeEnvRow(this)">×</button>
                </div>
            </div>
            <button id="add-env-btn">+ 添加变量</button>
        </div>
    </div>
	    <div id="result-modal" class="modal-overlay">
        <div class="modal">
            <h4 id="modal-title">运行结果</h4>
            <p id="modal-message"></p>
            <div class="modal-actions">
                <button class="modal-btn btn-close" id="modal-close-btn">关闭</button>
                <button class="modal-btn btn-view" id="modal-view-btn">查看运行结果</button>
            </div>
        </div>
    </div>
    <script>
        const runBtn = document.getElementById('run-btn');
        const addEnvBtn = document.getElementById('add-env-btn');
        const envList = document.getElementById('env-list');
        const codeInput = document.getElementById('code-input');
		const modal = document.getElementById('result-modal');
        const modalMsg = document.getElementById('modal-message');
        const modalCloseBtn = document.getElementById('modal-close-btn');
        const modalViewBtn = document.getElementById('modal-view-btn');
        function showModal(message, showViewBtn = true) {
            modalMsg.textContent = message;
            modalViewBtn.style.display = showViewBtn ? 'inline-block' : 'none';
            modal.style.display = 'flex';
        }
        function closeModal() {
            modal.style.display = 'none';
        }
        modalCloseBtn.addEventListener('click', closeModal);
        modalViewBtn.addEventListener('click', () => {
		modal.style.display = 'none';
		window.open('/test', '_blank');
        });
        runBtn.addEventListener('click', async () => {
            const code = codeInput.value;
            const envVars = getEnvVariables();
            let payload_obj = {
				payload_code: code
			};
			let savecode = {
				code: payload_obj,
				env: envVars
			};
  try {
    const response = await fetch("/dev_page", {
      method: "POST",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify(savecode)
    });
    const result = await response.json();
    console.log("result：", result);
    if (result.status == 'ok') {
        showModal('运行成功，是否查看运行结果？', true);
    } else {
        showModal('运行失败，请检查代码或环境变量。', false);
    }
  } catch (err) {
    console.log("result：", result);
	console.error("errcode：", err);
    showModal('上传失败，请检查网络或服务器状态。', false);
  }
        });
        addEnvBtn.addEventListener('click', () => {
            const row = document.createElement('div');
            row.className = 'env-row';
            row.innerHTML = `+"\x60"+`
                <input type="text" class="env-input key" placeholder="KEY">
                <input type="text" class="env-input val" placeholder="VALUE">
                <button class="icon-btn delete-btn" onclick="removeEnvRow(this)">×</button>
            `+"\x60"+`;
            envList.appendChild(row);
            row.querySelector('.key').focus();
        });
        window.removeEnvRow = function(btn) {
            const row = btn.parentElement;
            envList.removeChild(row);
        };
        function getEnvVariables() {
            const vars = {};
            const rows = document.querySelectorAll('.env-row');
            rows.forEach(row => {
                const key = row.querySelector('.key').value.trim();
                const val = row.querySelector('.val').value.trim();
                if (key) {
                    vars[key] = val;
                }
            });
            return vars;
        }
        codeInput.addEventListener('keydown', function(e) {
            if (e.key == 'Tab') {
                e.preventDefault();
                var start = this.selectionStart;
                var end = this.selectionEnd;
                this.value = this.value.substring(0, start) +
                    "  " + this.value.substring(end);
                this.selectionStart = this.selectionEnd = start + 2;
            }
        });
    </script>
</body>
</html>`