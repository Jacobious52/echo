window.onload = function() {
    var conn;
    var msg = document.getElementById("msg");
    var log = document.getElementById("log");

    var username = window.location.hash.substring(1);

    if (!username) {
        document.getElementById("form").setAttribute("hidden", "hidden");
        return
    }

    function appendLog(item) {
        var doScroll = log.scrollTop === log.scrollHeight - log.clientHeight;
        log.appendChild(item);
        if (doScroll) {
            log.scrollTop = log.scrollHeight - log.clientHeight;
        }
    }

    document.getElementById("form").onsubmit = function() {
        if (!conn) {
            return false;
        }
        if (!msg.value) {
            return false;
        }

        var jsonData = {};
        jsonData.Type = "msg";
        jsonData.User = "@" + username;
        jsonData.Msg = msg.value;

        conn.send(JSON.stringify(jsonData));

        msg.value = "";
        return false;
    };

    if (window["WebSocket"]) {
        conn = new WebSocket("ws://{{$}}/ws");

        conn.onclose = function(evt) {
            var item = document.createElement("div");
            item.innerHTML = "<b>Connection closed.</b>";
            appendLog(item);
        };

        conn.onopen = function(evt) {
            var jsonData = {};
            jsonData.Type = "conn";
            jsonData.User = "@" + username;
            jsonData.Msg = "has connected";

            conn.send(JSON.stringify(jsonData));
        }

        conn.onmessage = function(evt) {
            var messages = evt.data.split('\n\n');
            for (var i = 0; i < messages.length; i++) {

                var jsonData = JSON.parse(messages[i]);

                var item = document.createElement("div");

                switch (jsonData.Type) {
                    case "msg":
                        item.setAttribute("class", "alert alert-info");
                        break;
                    case "disconn":
                        item.setAttribute("class", "alert alert-danger");
                        break;
                    case "conn":
                        item.setAttribute("class", "alert alert-success");
                        break;
                    default:
                        console.error("unknown message type")
                }

                item.innerHTML = "<b>" + jsonData.User + "</b>: " + jsonData.Msg;
                appendLog(item);
            }
        };

    } else {
        var item = document.createElement("div");
        item.innerHTML = "<b>Your browser does not support WebSockets.</b>";
        appendLog(item);
    }
};
