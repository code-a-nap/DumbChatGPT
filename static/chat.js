const bot = "/static/bot.jpg";
const user = "/static/user.jpg";

function typeText(element, text){
    let index = 0;
  
    let interval = setInterval(() => {
      if(index < text.length){
        element.innerHTML += text.charAt(index);
        index++;
      }
      else
        clearInterval(interval);
    }, 30);
  }

const handleSubmit = async (e) => {
    e.preventDefault();
    var text = document.getElementById("content").value;
    document.getElementById("content").value = "";

    const chatContainer = document.querySelector('#chat_container');
    chatContainer.innerHTML += chatStripe(false, text);
    chatContainer.scrollTop = chatContainer.scrollHeight;

    await fetch('/chat/message', {
        method: 'POST',
        headers: {
        'Content-Type': 'application/json'
        },
        body: JSON.stringify({
        name: "user1",
        content: text
        })
    })
    .then((res) => res.json())
    .then((data) => {
        chatContainer.innerHTML += chatStripe(true, "");
        var messageElems = document.querySelectorAll(".message")
        typeText(messageElems[messageElems.length -1], data.content);
        chatContainer.scrollTop = chatContainer.scrollHeight;
    })
    .catch((err) => {
        console.log(err);
    });
};

function chatStripe(isAI, value){
    return (
      `
        <div class="wrapper ${isAI && 'ai'}">
          <div class="chat">
            <div class="profile">
              <img 
                src="${isAI ? bot : user}"
                alt="${isAI ? 'bot' : 'user'}"
              />
            </div>
            <div class="message">${value}</div>
            ${isAI ? '<button class="btn"><i class="fa fa-thumbs-up" style="color:#dcdcdc;"></i></button><button class="btn"><i class="fa fa-thumbs-down" style="color:#dcdcdc;"></i></button>' : ''}
          </div>
        </div>
      `
    )
  }

window.addEventListener('DOMContentLoaded', (event) => {
    
    var form = document.querySelector('form');
    form.addEventListener('submit', handleSubmit);
    form.addEventListener('keyup',async  (e) => {
        if (e.keyCode === 13)
          await handleSubmit(e);
      });
    
    var logoutBtn = document.getElementById("logout");
    logoutBtn.addEventListener("click", (e) => {
    document.cookie = "token=;expires=Thu, 01 Jan 1970 00:00:00 GMT; path=/;";
    window.location.href = "/";
    });

});
