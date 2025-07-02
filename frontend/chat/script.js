let socket;
let Port = 8080;

    function createRoom() {
      const name = document.getElementById("roomName").value;
      if (!name) return alert("กรุณากรอกชื่อห้อง");

      fetch(`http://localhost:${Port}/api/chat/create-room`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ name })
      })
      .then(res => res.json())
      .then(data => {
        if (data.room_id) {
          document.getElementById("roomId").value = data.room_id;
          document.getElementById("roomIdDisplay").textContent = "Room ID: " + data.room_id;
          alert("✅ สร้างห้องสำเร็จ!");
        } else {
          alert("❌ สร้างห้องไม่สำเร็จ");
        }
      });
    }

    function connectWebSocket() {
      const roomId = document.getElementById("roomId").value;
      const sender = document.getElementById("sender").value;

      if (!roomId || !sender) {
        return alert("กรุณากรอก Room ID และชื่อผู้ส่งก่อนเชื่อมต่อ");
      }

      socket = new WebSocket(`ws://localhost:${Port}/api/chat/ws?room_id=${roomId}&sender=${sender}`);

      socket.onopen = () => {
        appendChat("🟢 เชื่อมต่อห้องเรียบร้อยแล้ว", "system");
      };

      socket.onmessage = (event) => {
        appendChat(event.data, "system");
      };

      socket.onerror = (err) => {
        console.error("WebSocket error:", err);
      };
    }

    function sendMessage() {
      const msg = document.getElementById("message").value;
      const sender = document.getElementById("sender").value;
      if (!sender || !msg) {
        return alert("กรอกชื่อและข้อความก่อนส่ง");
      }

      if (!socket || socket.readyState !== WebSocket.OPEN) {
        return alert("กรุณาเชื่อมต่อห้องแชทก่อน");
      }

      socket.send(msg);
      appendChat(msg, sender);
      document.getElementById("message").value = "";
    }

    function appendChat(text, sender) {
      const chat = document.getElementById("chat");
      const div = document.createElement("div");
      div.className = "msg";
      div.innerHTML = `<span class="sender">${sender}:</span> ${text}`;
      chat.appendChild(div);
      chat.scrollTop = chat.scrollHeight;
    }