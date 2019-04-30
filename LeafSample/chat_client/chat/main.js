$(function() {
  var FADE_TIME = 150; // ms
  var TYPING_TIMER_LENGTH = 400; // ms
  var COLORS = [
    '#e21400', '#91580f', '#f8a700', '#f78b00',
    '#58dc00', '#287b00', '#a8f07a', '#4ae8c4',
    '#3b88eb', '#3824aa', '#a700ff', '#d300e7'
  ];

  // Initialize varibles
  var $window = $(window);
  var $usernameInput = $('.usernameInput'); // Input for username
  var $messages = $('.messages'); // Messages area
  var $inputMessage = $('.inputMessage'); // Input message input box

  var $loginPage = $('.login.page'); // The login page
  var $chatPage = $('.chat.page'); // The chatroom page

  // Prompt for setting a username
  var username;
  var connected = false;
  var typing = false;
  var lastTypingTime;
  var $currentInput = $usernameInput.focus();

  var socket = new WebSocket('ws://name5566.com:8888');
  var decoder = new TextDecoder('utf-8')
  socket.binaryType = 'arraybuffer';

  function socketSend (o, silence) {
    if (socket.readyState != socket.OPEN) {
      if (!silence) {
        addChatMessage({
          UserName: 'SYSTEM',
          Message: 'connection closed'
        });
      }
      return
    }
    socket.send(JSON.stringify(o))
  }

  function addParticipantsMessage (data) {
    var message = '';
    if (data.NumUsers === 1) {
      message += "there's 1 participant";
    } else {
      message += "there are " + data.NumUsers + " participants";
    }
    log(message);
  }

  // Sets the client's username
  function setUsername () {
    username = cleanInput($usernameInput.val().trim());

    // If the username is valid
    if (username) {
      $loginPage.fadeOut();
      $chatPage.show();
      $loginPage.off('click');
      $currentInput = $inputMessage.focus();

      // Tell the server your username
      socketSend({C2S_AddUser: {
        UserName: username
      }});
    }
  }

  // Sends a chat message
  function sendMessage () {
    var message = $inputMessage.val();
    // Prevent markup from being injected into the message
    message = cleanInput(message);
    // if there is a non-empty message and a socket connection
    if (message && connected) {
      $inputMessage.val('');
      addChatMessage({
        UserName: username,
        Message: message
      });
      // tell server to execute 'new message' and send along one parameter
      socketSend({C2S_Message: {
        Message: message
      }});
    }
  }

  // Log a message
  function log (message, options) {
    var $el = $('<li>').addClass('log').text(message);
    addMessageElement($el, options);
  }

  var imgReg = /:img\s+(\S+)/     
  // Adds the visual chat message to the message list
  function addChatMessage (data, options) {
    // Don't fade the message in if there is an 'X was typing'
    var $typingMessages = getTypingMessages(data);
    options = options || {};
    if ($typingMessages.length !== 0) {
      options.fade = false;
      $typingMessages.remove();
    }

    var $usernameDiv = $('<span class="username"/>')
      .text(data.UserName)
      .css('color', getUsernameColor(data.UserName));
    var regRes = imgReg.exec(data.Message)
    if (regRes != null) {
      var $messageBodyDiv = $('<img src="' + regRes[1] + '">');
    } else {
      var $messageBodyDiv = $('<span class="messageBody">')
        .text(data.Message);
    }

    var typingClass = data.Typing ? 'typing' : '';
    var $messageDiv = $('<li class="message"/>')
      .data('username', data.UserName)
      .addClass(typingClass)
      .append($usernameDiv, $messageBodyDiv);

    addMessageElement($messageDiv, options);
  }

  // Adds the visual chat typing message
  function addChatTyping (data) {
    data.Typing = true;
    data.Message = 'is typing';
    addChatMessage(data);
  }

  // Removes the visual chat typing message
  function removeChatTyping (data) {
    getTypingMessages(data).fadeOut(function () {
      $(this).remove();
    });
  }

  // Adds a message element to the messages and scrolls to the bottom
  // el - The element to add as a message
  // options.fade - If the element should fade-in (default = true)
  // options.prepend - If the element should prepend
  //   all other messages (default = false)
  function addMessageElement (el, options) {
    var $el = $(el);

    // Setup default options
    if (!options) {
      options = {};
    }
    if (typeof options.fade === 'undefined') {
      options.fade = true;
    }
    if (typeof options.prepend === 'undefined') {
      options.prepend = false;
    }

    // Apply options
    if (options.fade) {
      $el.hide().fadeIn(FADE_TIME);
    }
    if (options.prepend) {
      $messages.prepend($el);
    } else {
      $messages.append($el);
    }
    $messages[0].scrollTop = $messages[0].scrollHeight;
  }

  // Prevents input from having injected markup
  function cleanInput (input) {
    return $('<div/>').text(input).text();
  }

  // Updates the typing event
  function updateTyping () {
    if (connected) {
      if (!typing) {
        typing = true;
        socketSend({C2S_Action: {
          Op: 'typing'
        }}, true);
      }
      lastTypingTime = (new Date()).getTime();

      setTimeout(function () {
        var typingTimer = (new Date()).getTime();
        var timeDiff = typingTimer - lastTypingTime;
        if (timeDiff >= TYPING_TIMER_LENGTH && typing) {
          socketSend({C2S_Action: {
            Op: 'stop typing'
          }}, true);
          typing = false;
        }
      }, TYPING_TIMER_LENGTH);
    }
  }

  // Gets the 'X is typing' messages of a user
  function getTypingMessages (data) {
    return $('.typing.message').filter(function (i) {
      return $(this).data('username') === data.UserName;
    });
  }

  // Gets the color of a username through our hash function
  function getUsernameColor (username) {
    // Compute hash code
    var hash = 7;
    for (var i = 0; i < username.length; i++) {
       hash = username.charCodeAt(i) + (hash << 5) - hash;
    }
    // Calculate color
    var index = Math.abs(hash % COLORS.length);
    return COLORS[index];
  }

  // Keyboard events

  $window.keydown(function (event) {
    // Auto-focus the current input when a key is typed
    if (!(event.ctrlKey || event.metaKey || event.altKey)) {
      $currentInput.focus();
    }
    // When the client hits ENTER on their keyboard
    if (event.which === 13) {
      if (username) {
        sendMessage();
        socketSend({C2S_Action: {
          Op: 'stop typing'
        }}, true);
        typing = false;
      } else {
        setUsername();
      }
    }
  });

  $inputMessage.on('input', function() {
    updateTyping();
  });

  // Click events

  // Focus input when clicking anywhere on login page
  $loginPage.click(function () {
    $currentInput.focus();
  });

  // Focus input when clicking on the message input's border
  $inputMessage.click(function () {
    $inputMessage.focus();
  });

  // Socket events
  socket.onmessage = function(e) {
    var data = JSON.parse(decoder.decode(e.data));

    // Whenever the server emits 'login', log the login message
    if (data.S2C_Login) {
      connected = true;
      // Display the welcome message
      var message = "Leaf game framework: https://github.com/name5566/leaf";
      log(message, {
        prepend: true
      });
      addParticipantsMessage(data.S2C_Login);
    }

    // Whenever the server emits 'user joined', log it in the chat body
    if (data.S2C_Joined) {
      log(data.S2C_Joined.UserName + ' joined');
      addParticipantsMessage(data.S2C_Joined);
    }

    // Whenever the server emits 'user left', log it in the chat body
    if (data.S2C_Left) {
      log(data.S2C_Left.UserName + ' left');
      addParticipantsMessage(data.S2C_Left);
      removeChatTyping(data.S2C_Left);
    }

    // Whenever the server emits 'typing', show the typing message
    if (data.S2C_Typing) {
      addChatTyping(data.S2C_Typing);
    }

    // Whenever the server emits 'stop typing', kill the typing message
    if (data.S2C_StopTyping) {
      removeChatTyping(data.S2C_StopTyping);
    }

    // Whenever the server emits 'new message', update the chat body
    if (data.S2C_Message) {
      addChatMessage(data.S2C_Message);
    }
  }

  socket.onclose = function() {
    addChatMessage({
      UserName: 'SYSTEM',
      Message: 'connection closed'
    });
  }

  socket.onerror = function() {
    addChatMessage({
      UserName: 'SYSTEM',
      Message: 'connection closed'
    });
  }
});
