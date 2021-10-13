function d(a) {
  return document.getElementById(a);
}

function s() {
  socket.emit("SetName", d("me").value);
  ls("me", d("me").value);
}

function m() {
  if (event.keyCode == 13 && d("msg").value) {
    Dream(d("to").value + "|" + d("msg").value);
    d("msgs").innerHTML =
      d("me").value + ": " + d("msg").value + "\n" + d("msgs").innerHTML;
    d("msg").value = "";
    ls("msgs", d("msgs").innerHTML);
  }
}

function Dream(data) {
  if (typeof data === "object") {
    data = JSON.stringify(data);
  }
  socket.emit("Dream", data);
}

function Debounce(func, wait, immediate) {
  var timeout;
  return function () {
    var context = this,
      args = arguments;
    var later = function () {
      timeout = null;
      if (!immediate) func.apply(context, args);
    };
    var callNow = immediate && !timeout;
    clearTimeout(timeout);
    timeout = setTimeout(later, wait);
    if (callNow) func.apply(context, args);
  };
}

function Bildirim(title, text, img) {
  Notification.requestPermission().then(function (result) {
    /*console.log(result);*/
  });
  if (document.hidden) {
    new Notification(title, { body: text, icon: img });
  }
}

function ls(key, obj) {
  key = window.top.location.href + "_" + key;
  if (obj) {
    localStorage.setItem(key, JSON.stringify(obj));
  }
  if (localStorage.getItem(key) != null) {
    return JSON.parse(localStorage.getItem(key));
  } else {
    return [];
  }
}

function clearMsgs() {
  ls("msgs", []);
  d("msgs").innerHTML = "";
}

window.addEventListener("load", function () {
  window.socket = io();

  socket.on("connect", function () {
    console.log("connected");
    setTimeout(s, 100);
  });

  socket.on("Dream", function (a) {
    a = a.split("|");
    if (a.length > 1) {
      var msg = a[0] + ": " + a[1];

      Bildirim(
        "BackMessage",
        msg,
        "https://yt3.ggpht.com/ytc/AKedOLQzo8YfpM98YDjaEwUPL-mM1XctcrOMl4DuggZoxA=s100-c-k-c0x00ffffff-no-rj"
      );

      d("msgs").innerHTML = msg + "\n" + d("msgs").innerHTML;

      ls("msgs", d("msgs").innerHTML);
    }
  });

  Notification.requestPermission().then(function (result) {
    /*console.log(result);*/
  });

  d("me").value = ls("me");
  s();

  d("msgs").innerHTML = ls("msgs");

  d("me").addEventListener(
    "keyup",
    Debounce(function () {
      s();
    }, 500)
  );
});
