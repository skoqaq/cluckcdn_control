function getIpFlag(ip){
    let url;
    $.ajaxSettings.async = false;
    $.getJSON("https://api.ip.sb/geoip/"+ip,
        function(json) {
            if (json['country_code'] != undefined) {
                //url = `https://www.countryflags.io/${json['country_code']}/flat/64.png`;
                url_raw = `https://flagcdn.com/w40/${json['country_code'].png}`;
                console.log(url_raw.toLowerCase());
                url_new = url_raw.toLowerCase();
                url = `https://flagcdn.com/w40/${["url_new"]}.png`;
            }else{
                url = `https://www.countryflags.io/hk/flat/64.png`;
            }
            
        }
    );
    return url;
}
let status, nodeList;
function getNode(s){
    $.ajax({
        url: "/api/nodeList",
        async: true,
        type: "POST",
        success: function(re){
            let nodelist = document.getElementById('nodelist');
            nodelist.innerHTML = "";
            for (let i in re){
                let color, text;
                if (re[i]["status"] == "false"){
                    color = "orange";
                    text = `<div class="info-text"><b>IP: </b>${re[i]["ip"]}</div><div class="info-text"><b>System: </b>${re[i]["system"]}</div>`;
                }else if (re[i]["status"] == "true"){
                    color = "green";
                    text = `<div class="info-text"><b>IP: </b>${re[i]["ip"]}</div><div class="info-text"><b>System: </b>${re[i]["system"]}</div>`;
                }else if (re[i]["status"] == "error"){
                    color = "red";
                    text = `<div class="info-text"><b>IP: </b>${re[i]["ip"]}</div><div class="info-text"><b>Error: </b>Unable to connect </div>`;
                }
                nodelist.innerHTML += `<div class="mdui-col">
    <div class="mdui-card card-m">
        <div class="mdui-m-a-2 card-info">
            <img style="float: left;" src="${getIpFlag(re[i]["ip"])}">
            <div style="float: left;" class="l-text">
                <span style="font-size: 120%;font-weight: 500;">${re[i]["name"]}</span>
                <div style="margin-top: 5px;">
                    ${text}
                </div>
            </div>
            <div style="float: right;" class="l-button">
                <i class="mdui-icon material-icons" style="cursor: pointer;" mdui-dialog="{target: '#${re[i]["name"]}-set'}">settings</i>
                <i class="mdui-icon material-icons" data-ip="${re[i]["ip"]}" name="status" style="color: ${color};">fiber_manual_record</i>
            </div>
        </div>
    </div>
</div>
<div class="mdui-dialog" id="${re[i]["name"]}-set">
    <div class="mdui-dialog-title">${re[i]["name"]}</div>
    <div class="mdui-dialog-content" style="padding: 0 24px;">
        <div class="mdui-btn-group" style="width: 100%;">
            <button type="button" class="mdui-btn cluck-button" onclick="setNodeStatus('${re[i]["ip"]}', 'setNodeStatus', 'start')"><i class="mdui-icon material-icons">power_settings_new</i> Start</button>
            <button type="button" class="mdui-btn cluck-button" onclick="setNodeStatus('${re[i]["ip"]}', 'setNodeStatus', 'stop')"><i class="mdui-icon material-icons">power_settings_new</i> Stop</button>
            <button type="button" class="mdui-btn cluck-button" onclick="setNodeStatus('${re[i]["ip"]}', 'setNodeStatus', 'reload')"><i class="mdui-icon material-icons">autorenew</i> ReStart</button>
            <button type="button" class="mdui-btn cluck-button" onclick="setNodeStatus('${re[i]["ip"]}', 'updateVhost')"><i class="mdui-icon material-icons">cloud_upload</i> Update</button>
        </div>
        <button type="button" class="mdui-btn mdui-btn-block" mdui-dialog="{target: '#${re[i]["name"]}-del'}" mdui-dialog-close><i class="mdui-icon material-icons">delete</i> Delete</button>
    </div>
    <div class="mdui-dialog-actions">
        <button class="mdui-btn mdui-ripple" mdui-dialog-close>cancel</button>
    </div>
</div>
<div class="mdui-dialog" id="${re[i]["name"]}-del">
    <div class="mdui-dialog-title">Delete ${re[i]["name"]}</div>
    <div class="mdui-dialog-content">
        你真的要刪掉這個節點嗎？
    </div>
    <div class="mdui-dialog-actions">
        <button class="mdui-btn mdui-ripple" mdui-dialog-close>不要！</button>
        <button class="mdui-btn mdui-ripple" onclick="delNode('${re[i]["ip"]}')" mdui-dialog-close>刪掉！</button>
    </div>
</div>`; 
            }
            // Update
            status = document.getElementsByName("status");
            nodeList = [];
            for (let i=0;i<status.length;i++){
                nodeList.push(status[i].getAttribute("data-ip"));
            }
            if (s != "reload"){
                setTimeout(function(){updateAsync()}, 10000);
            }
        },
        error: function(){
            mdui.snackbar({
                message: 'API Connect Error',
                position: 'left-bottom'
            });
        }
    })
}
function updateAsync(){
    $.ajax({
        url: "/api/updateNodeStatus",
        type: "POST",
        async: true,
        data: { "nodelist": JSON.stringify(nodeList) },
        success: function (re) {
            for (let i = 0; i < status.length; i++) {
                if (re[i] == "true") {
                    status[i].style.color = "green";
                } else if (re[i] == "false") {
                    status[i].style.color = "orange";
                } else if (re[i] == "error") {
                    status[i].style.color = "red";
                }
            }
            setTimeout(function(){updateAsync(status, nodeList)}, 10000);
        },
        error: function () {
            mdui.snackbar({
                message: 'API Connect Error',
                position: 'left-bottom'
            });
        }
    });
}
function snackbar(status){
    switch(status){
        case "startOk":
            mdui.snackbar({
                message: 'Start Node',
                position: 'right-bottom'
            });
            break;
        case "isStarting":
            mdui.snackbar({
                message: '節點正在運行',
                position: 'right-bottom'
            });
            break;
        case "stopOk":
            mdui.snackbar({
                message: 'Stop Node',
                position: 'right-bottom'
            });
            break;
        case "isStopped":
            mdui.snackbar({
                message: '節點已經閂咗',
                position: 'right-bottom'
            });
            break;
        case "reloadOk":
            mdui.snackbar({
                message: '節點重載',
                position: 'right-bottom'
            });
            break;
        case "upOk":
            mdui.snackbar({
                message: '節點同步',
                position: 'right-bottom'
            });
            break;
        case "updateOk":
            mdui.snackbar({
                message: '節點同步',
                position: 'right-bottom'
            });
            break;
    }
}
function setNodeStatus(ip, api, status){
    $.ajax({
        url: "/api/postNodeApi",
        type: "POST",
        async: true,
        data: {"ip": ip, "api": api, "status": status},
        success: function(re){
            snackbar(re['status']);
        },
        error: function(){
            mdui.snackbar({
                message: 'API Connect Error',
                position: 'left-bottom'
            });
        }
    })
}
function setAllNodeStatus(status){
    $.ajax({
        url: "/api/setAllStatus",
        type: "POST",
        async: true,
        data: {"status": status},
        success: function(re){
            snackbar(re['status']);
        },
        error: function(){
            mdui.snackbar({
                message: 'API Connect Error',
                position: 'left-bottom'
            });
        }
    })
}
function delNode(ip){
    $.ajax({
        url: "/api/delNode",
        type: "POST",
        async: true,
        data: {"ip": ip},
        success: function(re){
            if (re["status"] == "ok"){
                getNode("reload");
            }else{
                mdui.snackbar({
                    message: 'API Error',
                    position: 'left-bottom'
                });
            }
        },
        error: function(){
            mdui.snackbar({
                message: 'API Connect Error',
                position: 'left-bottom'
            });
        }
    })
}

function addNode(){
    let name = document.getElementById("add_name").value;
    let ip = document.getElementById("add_ip").value;
    if (name != "" && ip != ""){
        $.ajax({
            url: "/api/addNode",
            type: "POST",
            async: true,
            data: {"name": name, "ip": ip},
            success: function(re){
                if (re["status"] == "ok"){
                    mdui.snackbar({
                        message: 'Add OK',
                        position: 'left-bottom'
                    });
                    getNode("reload");
                }else if (re["status"] == "repeat"){
                    mdui.snackbar({
                        message: '已經添加過',
                        position: 'left-bottom'
                    });
                }
            },
            error: function(){
                mdui.snackbar({
                    message: 'API Connect Error',
                    position: 'left-bottom'
                });
            }
        })
    }else{
        mdui.snackbar({
            message: '未填寫 Form',
            position: 'left-bottom'
        });
    }
}

function delWebsite(host){
    $.ajax({
        url: "/api/delWebSite",
        type: "POST",
        async: true,
        data: {"host": host},
        success: function(re){
            if (re["status"] == "ok"){
                getWebSite();
            }else{
                mdui.snackbar({
                    message: 'API Error',
                    position: 'left-bottom'
                });
            }
        },
        error: function(){
            mdui.snackbar({
                message: 'API Connect Error',
                position: 'left-bottom'
            });
        }
    })
}

function https(){
    let use_https = document.getElementById("use_https").checked;
    let pen = document.getElementById('pen');
    let key = document.getElementById('key');
    if (!use_https){
        pen.setAttribute("disabled", "");
        key.setAttribute("disabled", "");
        pen.value = "";
        key.value = "";
    }else{
        pen.removeAttribute("disabled");
        key.removeAttribute("disabled");
    }
}
function addWebSite(){
    let name = document.getElementById("add_name");
    let host = document.getElementById("add_host");
    let source_host = document.getElementById("add_source_host");
    let source = document.getElementById("add_source");
    let text = document.getElementById("add_text");
    let proto = document.getElementsByName("proto");
    for (let i in proto){
        if (proto[i].checked == true){
            proto = proto[i].getAttribute("data-tls");
            break;
        }
    }
    let pen = document.getElementById("pen");
    let key = document.getElementById("key");

    $.ajax({
        url: "/api/addWebSite",
        type: "POST",
        async: true,
        data: {"name": name.value, "host": host.value, "source_host": source_host.value, "source": source.value, "text": text.value, "proto": proto, "pen": pen.value, "key": key.value},
        success: function(re){
            if (re["status"] == "ok"){
                name.value = "";
                host.value = "";
                source_host.value = "";
                source.value = "";
                text.value = "";
                pen.value = "";
                key.value = "";
                getWebSite();
            }else{
                mdui.snackbar({
                    message: '請填寫完整',
                    position: 'left-bottom'
                });
            }
        },
        error: function(){
            mdui.snackbar({
                message: 'API Connect Error',
                position: 'left-bottom'
            });
        }
    })
}
function getWebSite(){
    $.ajax({
        url: "/api/websiteList",
        type: "POST",
        async: true,
        success: function(re){
            let nodelist = document.getElementById('nodelist');
            nodelist.innerHTML = "";
            for (let i in re){
                nodelist.innerHTML += `<div class="mdui-col">
    <div class="mdui-card card-m">
        <div class="mdui-m-a-2 card-info" style="height: 160px;">
            <div style="float: left;">
                <span style="font-size: 130%;font-weight: 500;">${re[i]["name"]}</span>
                <div style="margin-top: 5px;line-height: 21px;">
                    <span class="web-2">Host: </span><span>${re[i]["host"]}</span><br>
                    <span class="web-2">源Host: </span><span>${re[i]["source_host"]}</span><br>
                    <span class="web-2">源伺服器: </span><span>${re[i]["source"]}</span><br>
                    <span class="web-2">回源協議: </span><span>${re[i]["proto"]}</span><br>
                    <span class="web-2">使用TLS: </span><span>${re[i]["tls"]}</span><br>
                    <span class="web-2">網站註解: </span><span>${re[i]["text"]}</span>
                </div>
            </div>
            <div style="float: right;" class="l-button">
                <i class="mdui-icon material-icons" style="cursor: pointer;" mdui-dialog="{target: '#${re[i]["name"]}-del'}">delete</i>
            </div>
        </div>
    </div>
</div>
<div class="mdui-dialog" id="${re[i]["name"]}-del">
    <div class="mdui-dialog-title">Delete ${re[i]["name"]}</div>
    <div class="mdui-dialog-content">
        你真的要刪掉這個網站嗎？
    </div>
    <div class="mdui-dialog-actions">
        <button class="mdui-btn mdui-ripple" mdui-dialog-close>不要！</button>
        <button class="mdui-btn mdui-ripple" onclick="delWebsite('${re[i]["host"]}')" mdui-dialog-close>刪掉！</button>
    </div>
</div>`;
            }
        },
        error: function(){
            mdui.snackbar({
                message: 'API Connect Error',
                position: 'left-bottom'
            });
        }
    })
}
