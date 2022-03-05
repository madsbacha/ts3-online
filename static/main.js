let ws;
const onlineEl = document.getElementById("online");
const usersEl = document.getElementById("users");

function updateStatus(online, usernames) {
    // Update online count
    onlineEl.innerHTML = online;

    // Remove all existing li elements
    while (usersEl.firstChild) {
        usersEl.removeChild(usersEl.lastChild);
    }

    // Add new li elements with new users
    for(let i = 0; i < usernames.length; i++){
        const liEl = document.createElement("li");
        liEl.innerHTML = usernames[i];
        usersEl.appendChild(liEl);
    }
}

async function onMessage(evt) {
    const res = await evt.data.text()
    var resJson = JSON.parse(res);

    updateStatus(resJson.online, resJson.usernames)
}

function onClose(e) {
    console.log('Socket is closed. Reconnect will be attempted in 1 second.', e.reason);
    setTimeout(function() {
        connect();
    }, 1000);
}

function onError(e) {
    console.log("ERROR: ", e);
    ws.close();
}

function connect() {
    const protocol = location.protocol.startsWith("https") ? "wss" : "ws"
    ws = new WebSocket(`${protocol}://${location.host}/websocket`);
    ws.onclose = onClose;
    ws.onmessage = onMessage;
    ws.onerror = onError;
}

connect()
