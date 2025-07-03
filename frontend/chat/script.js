
let BaseURL = "https://renewed-similarly-collie.ngrok-free.app";
let socket;

/**
 * @param {string} message
 * @param {string} type
 */
function showMessageBox(message, type = "info") {
    const messageBox = document.getElementById("messageBox");
    const messageText = document.getElementById("messageText");
    messageText.textContent = message;

    messageBox.style.backgroundColor = "";
    if (type === "success") {
        messageBox.style.backgroundColor = "#d4edda";
        messageText.style.color = "#155724";
    } else if (type === "error") {
        messageBox.style.backgroundColor = "#f8d7da";
        messageText.style.color = "#721c24";
    } else {
        messageBox.style.backgroundColor = "#cce5ff";
        messageText.style.color = "#004085";
    }

    messageBox.style.display = "block";
    setTimeout(() => {
        messageBox.style.display = "none";
    }, 3000);
}

/**
 * @param {string} message
 * @param {string} type
 * @param {string} sender
 */
function appendChat(message, type, sender = "") {
    const chatDiv = document.getElementById("chat");
    const msgElement = document.createElement("div");
    msgElement.classList.add("msg");

    if (type === "system") {
        msgElement.style.color = "#888";
        msgElement.innerHTML = `<em>${message}</em>`;
    } else {
        msgElement.innerHTML = `<span class="sender">${sender}:</span> ${message}`;
    }
    chatDiv.appendChild(msgElement);
    chatDiv.scrollTop = chatDiv.scrollHeight;
}

function createRoom() {
    const name = document.getElementById("roomName").value;
    if (!name) {
        showMessageBox("กรุณากรอกชื่อห้อง", "error");
        return;
    }

    fetch(`${BaseURL}/api/chat_rooms/create_room`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ name })
    })
    .then(res => res.json())
    .then(data => {
        if (data.room_id) {
            document.getElementById("roomId").value = data.room_id;
            document.getElementById("roomIdDisplay").textContent = "Room ID: " + data.room_id;
            showMessageBox("✅ สร้างห้องสำเร็จ!", "success");
        } else {
            showMessageBox("❌ สร้างห้องไม่สำเร็จ", "error");
        }
    })
    .catch(error => {
        console.error("Error creating room:", error);
        showMessageBox("❌ เกิดข้อผิดพลาดในการสร้างห้อง", "error");
    });
}

function connectWebSocket() {
  const roomId = document.getElementById("roomId").value;
  const sender = document.getElementById("sender").value;

  if (!roomId || !sender) {
    showMessageBox("กรุณากรอก Room ID และชื่อผู้ส่งก่อนเชื่อมต่อ", "error");
    return;
  }

  console.log("Fetch URL:", `${BaseURL}/api/chat_rooms/history?room_id=${roomId}`);

  fetch(`${BaseURL}/api/chat_rooms/history?room_id=${roomId}`, {
    method: "GET",
    headers: { "Content-Type": "application/json","ngrok-skip-browser-warning": "true"}
  })
  .then(res => res.json())
  .then(data => {
    console.log("Chat history data:", data);
    if (data.chat_messages) {
      data.chat_messages.forEach(msg => {
        appendChat(msg.message, "other", msg.sender);
      });
    }
  })
  .catch(error => {
    console.error("Error fetching chat history:", error);
    showMessageBox("❌ เกิดข้อผิดพลาดในการดึงประวัติแชท", "error");
  })
  .finally(() => {
    openWebSocket(roomId, sender);
  });
}



function openWebSocket(roomId, sender) {
    if (socket && socket.readyState === WebSocket.OPEN) {
        socket.close();
        appendChat("🟡 ตัดการเชื่อมต่อเก่าแล้ว", "system");
    }

    socket = new WebSocket(`wss://${BaseURL.replace("https://", "")}/api/chat_rooms/ws?room_id=${roomId}&sender=${sender}`);

    socket.onopen = () => {
        appendChat("🟢 เชื่อมต่อห้องเรียบร้อยแล้ว", "system");
        showMessageBox("🟢 เชื่อมต่อ WebSocket สำเร็จ!", "success");
    };

    socket.onmessage = (event) => {
        const [sender, message] = event.data.split("|");
        appendChat(message, "other", sender);
    };

    socket.onerror = (err) => {
        console.error("WebSocket error:", err);
        showMessageBox("❌ เกิดข้อผิดพลาดในการเชื่อมต่อ WebSocket", "error");
    };

    socket.onclose = (event) => {
        appendChat(`🔴 การเชื่อมต่อ WebSocket ปิดลง: ${event.code} ${event.reason}`, "system");
        if (!event.wasClean) {
            showMessageBox("🔴 การเชื่อมต่อ WebSocket ขาดหายไป", "error");
        }
    };
}


function sendMessage() {
    const messageInput = document.getElementById("message");
    const message = messageInput.value;
    const sender = document.getElementById("sender").value;

    if (!message || !socket || socket.readyState !== WebSocket.OPEN) {
        showMessageBox("ไม่สามารถส่งข้อความได้. ตรวจสอบการเชื่อมต่อและข้อความของคุณ", "error");
        return;
    }

    socket.send(message);
    appendChat(message, "user", sender);
    messageInput.value = "";
}