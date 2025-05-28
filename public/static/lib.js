/**
 * 设置cookie
 * @param {string} name cookie名称
 * @param {string} value cookie值
 * @param {number} exp 过期时间
 */
function setCookie(name, value, exp) {
  let expires = "";
  if (exp) {
    const date = new Date();
    if (exp === -1) {
      date.setTime(date.getTime() + exp * 24 * 60 * 60 * 1000);
    } else {
      date.setTime(exp * 1000);
    }
    expires = "; expires=" + date.toUTCString();
  }
  let cookie = name + "=" + encodeURIComponent(value) + expires + "; path=/";
  document.cookie = cookie;
}

/**
 * 获取cookie
 * @param {string} name cookie名称
 * @returns {string|null} cookie值，不存在时返回null
 */
function getCookie(name) {
  const nameEQ = name + "=";
  const ca = document.cookie.split(";");
  for (let i = 0; i < ca.length; i++) {
    let c = ca[i];
    while (c.charAt(0) === " ") {
      c = c.substring(1);
    }
    if (c.indexOf(nameEQ) === 0) {
      return decodeURIComponent(c.substring(nameEQ.length));
    }
  }
  return null;
}

/**
 * 删除cookie
 * @param {string} name cookie名称
 */
function deleteCookie(name) {
  setCookie(name, "", -1);
}

var msgTimeout;
function showMsg(content) {
  msgTimeout && clearTimeout(msgTimeout);
  const msgModal = document.getElementById("msgModal");
  if (msgModal.classList.contains("hidden")) {
    msgModal.classList.remove("hidden");
  }
  const msgContent = document.getElementById("msgContent");
  if (content !== undefined) {
    msgContent.innerHTML = content;
  } else {
    msgContent.innerHTML = "发生未知错误";
  }
  msgTimeout = setTimeout(() => {
    msgModal.classList.add("hidden");
  }, 3000);
}

function initModal(type, submitFunc) {
  const btn = document.getElementById(type+"Button");
  const modal = document.getElementById(type+"Modal");
  const closeModal = document.getElementById(type+"Close");
  const form = document.getElementById(type+"Form");
  btn.addEventListener("click", () => {
    modal.classList.remove("hidden");
  });
  closeModal.addEventListener("click", () => {
    modal.classList.add("hidden");
  });
  modal.addEventListener("click", (e) => {
    if (e.target === modal) {
      modal.classList.add("hidden");
    }
  });
  form.addEventListener("submit", (e) => {
    e.preventDefault();

    const button = form.querySelector('button[type="submit"]');
    button.disabled = true;
    button.classList.add("cursor-not-allowed");
    let lastInner = button.innerHTML;
    button.innerHTML = '<i class="ri-loop-left-line"></i>'+lastInner;

    submitFunc().then(() => {
      button.disabled = false;
      button.classList.remove("cursor-not-allowed");
      button.innerHTML = lastInner;
    });
  });
  const cancel = document.getElementById(type+"Cancel");
  cancel.addEventListener("click", (e) => {
    e.preventDefault();
    modal.classList.add("hidden");
  });
}

function closedModal(eleId) {
  return new Promise((resolve) => {
    const modal = document.getElementById(eleId);
    modal.classList.add("hidden");
    resolve();
  });
}

function refresh() {
  window.location.reload();
}

function showHome() {
  document.getElementById("defualt_empty").classList.add("hidden");
  document.getElementById("home").classList.remove("hidden");
  document.getElementById("content").classList.add("hidden");
  document.getElementById("toc").innerHTML = "";
  document.getElementById("markdownContent").innerHTML = "";
  document.getElementById("tocToggle").classList.add("hidden");
  document.getElementById("searchElement").classList.add("hidden");
  document.getElementById("titleBox").classList.remove("hidden");
  document.getElementById("miniTitleBox").classList.add("hidden");
}

function showContent() {
  document.getElementById("defualt_empty").classList.add("hidden");
  document.getElementById("home").classList.add("hidden");
  document.getElementById("home").innerHTML = "";
  document.getElementById("content").classList.remove("hidden");
  document.getElementById("tocToggle").classList.remove("hidden");
  document.getElementById("searchElement").classList.remove("hidden");
  document.getElementById("titleBox").classList.add("hidden");
  document.getElementById("miniTitleBox").classList.remove("hidden");
}

function showDefault() {
  document.getElementById("defualt_empty").classList.remove("hidden");
  document.getElementById("home").classList.add("hidden");
  document.getElementById("home").innerHTML = "";
  document.getElementById("content").classList.add("hidden");
  document.getElementById("toc").innerHTML = "";
  document.getElementById("markdownContent").innerHTML = "";
  document.getElementById("tocToggle").classList.add("hidden");
  document.getElementById("searchElement").classList.add("hidden");
  document.getElementById("titleBox").classList.remove("hidden");
  document.getElementById("miniTitleBox").classList.add("hidden");
}

function isLogin() {
  return getCookie("token") !== null;
}