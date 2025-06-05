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
  const btn = document.getElementById(type + "Button");
  const modal = document.getElementById(type + "Modal");
  const closeModal = document.getElementById(type + "Close");
  const form = document.getElementById(type + "Form");
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
    button.innerHTML =
      '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" width="16" height="16" fill="currentColor"><path fill="none" d="M0 0h24v24H0z"></path><path d="M12 4C9.25144 4 6.82508 5.38626 5.38443 7.5H8V9.5H2V3.5H4V5.99936C5.82381 3.57166 8.72764 2 12 2C17.5228 2 22 6.47715 22 12H20C20 7.58172 16.4183 4 12 4ZM4 12C4 16.4183 7.58172 20 12 20C14.7486 20 17.1749 18.6137 18.6156 16.5H16V14.5H22V20.5H20V18.0006C18.1762 20.4283 15.2724 22 12 22C6.47715 22 2 17.5228 2 12H4Z"></path></svg>' +
      lastInner;

    submitFunc().then(() => {
      button.disabled = false;
      button.classList.remove("cursor-not-allowed");
      button.innerHTML = lastInner;
    });
  });
  const cancel = document.getElementById(type + "Cancel");
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

function titleAnimate(show) {
  if (show === undefined) {
    show = false;
  }

  const titleBox = document.getElementById("titleBox");
  titleBox.className = show
    ? "transition-all duration-300 ease-in-out flex items-center justify-center flex-col space-y-6 mb-8"
    : "transition-all duration-300 ease-in-out flex items-center justify-center";

  const webTitle = document.getElementById("webTitle");
  webTitle.className = show
    ? "transition-all duration-300 ease-in-out text-4xl font-bold text-black dark:text-gray-400 text-center"
    : "transition-all duration-300 ease-in-out absolute top-[2.15rem] text-xl font-bold text-black dark:text-gray-400";

  const webDesc = document.getElementById("webDesc");
  const webNacFloat = document.getElementById("webNavFloat");
  webNacFloat.classList.remove("w-full");

  if (show) {
    webDesc.classList.remove("hidden");
    webNacFloat.classList.remove("absolute");
  } else {
    webDesc.classList.add("hidden");
    webNacFloat.classList.add("absolute");
  }

  const webNav = document.getElementById("webNav");
  webNav.className = show
    ? "flex justify-center items-center space-x-4 flex-wrap"
    : "flex flex-col space-x-2 text-xs leading-[0.75]";

  const navLinks = document.querySelectorAll(".nav-link");
  navLinks.forEach((link) => {
    link.className = show
      ? "nav-link text-nowrap sm:text-base text-sm text-gray-600 hover:text-primary dark:text-gray-300 dark:hover:text-gray-100"
      : "nav-link bg-blue-900 opacity-80 pl-1.5 pr-1 text-gray-200 hover:bg-blue-700 hover:text-gray-100 text-vertical rounded-l-sm px-2";
  });

  const navDelimiter = document.querySelectorAll(".nav-delimiter");
  navDelimiter.forEach((element) => {
    if (show) {
      element.classList.remove("hidden");
    } else {
      element.classList.add("hidden");
    }
  });
}

function showHome() {
  document.getElementById("loading").classList.add("hidden");
  document.getElementById("defualt_empty").classList.add("hidden");
  document.getElementById("home").classList.remove("hidden");
  document.getElementById("content").classList.add("hidden");
  document.getElementById("toc").innerHTML = "";
  document.getElementById("markdownContent").innerHTML = "";
  document.getElementById("tocToggle").classList.add("hidden");
  document.getElementById("searchElement").classList.remove("hidden");
}

function showContent() {
  document.getElementById("loading").classList.add("hidden");
  document.getElementById("defualt_empty").classList.add("hidden");
  document.getElementById("home").classList.add("hidden");
  document.getElementById("home").innerHTML = "";
  document.getElementById("content").classList.remove("hidden");
  document.getElementById("tocToggle").classList.remove("hidden");
  document.getElementById("searchElement").classList.remove("hidden");
}

function showTitleLoading() {
  // 初始化标题加载占位符
  document.getElementById("webTitle").innerHTML = buildTpl(loadingTpl, {
    height: "8",
    width: "1/4",
  });
  document.getElementById("webDesc").innerHTML = buildTpl(loadingTpl, {
    height: "4",
    width: "1/4",
  });
  document.getElementById("webNav").innerHTML = buildTpl(loadingTpl, {
    height: "4",
    width: "1/3",
  });
}

function showLoadding() {
  document.getElementById("defualt_empty").classList.add("hidden");
  document.getElementById("home").classList.add("hidden");
  document.getElementById("home").innerHTML = "";
  document.getElementById("content").classList.add("hidden");
  document.getElementById("toc").innerHTML = "";
  document.getElementById("markdownContent").innerHTML = "";
  document.getElementById("tocToggle").classList.add("hidden");
  document.getElementById("searchElement").classList.remove("hidden");

  // 初始化正文加载占位符
  const loading = document.getElementById("loading");
  // 生成3-6个随机大小的loading占位符
  let loadingHtml = "";
  const count = Math.floor(Math.random() * 4) + 3; // 3-6之间的随机数
  for (let i = 0; i < count; i++) {
    const height = Math.floor(Math.random() * 4) + 4; // 4-7之间的随机高度
    const widthOptions = ["1/4", "1/3", "1/2", "2/3", "3/4", "full"];
    const width = widthOptions[Math.floor(Math.random() * widthOptions.length)];
    loadingHtml +=
      buildTpl(loadingTpl, { height: height.toString(), width: width }) +
      '<div class="h-4"></div>';
  }
  loading.innerHTML = loadingHtml;
  loading.classList.remove("hidden");
}

function showDefault() {
  document.getElementById("loading").classList.add("hidden");
  document.getElementById("defualt_empty").classList.remove("hidden");
  document.getElementById("home").classList.add("hidden");
  document.getElementById("home").innerHTML = "";
  document.getElementById("content").classList.add("hidden");
  document.getElementById("toc").innerHTML = "";
  document.getElementById("markdownContent").innerHTML = "";
  document.getElementById("tocToggle").classList.add("hidden");
  document.getElementById("searchElement").classList.add("hidden");
  titleAnimate(true);
}

function isLogin() {
  return getCookie("token") !== null;
}

function buildTpl(tpl, data) {
  return tpl.replace(/{{\s*(\w+)\s*}}/g, (match, key) => {
    return data[key] !== undefined ? data[key] : "";
  });
}

function isHome() {
  const hash = window.location.hash.slice(1);
  return hash === "" || hash === "home" || hash === "index";
}

function cleanSearch() {
  document.getElementById("searchInput").value = "";
  document.getElementById("searchInput").focus();
  document.getElementById("searchInput").blur();
}

function highlightContent(elementId, searchTerm) {
  const container = document.getElementById(elementId);

  // 移除所有高亮
  const highlightRegex = /<span class="search-highlight">(.*?)<\/span>/g;
  const content = container.innerHTML;
  container.innerHTML = content.replace(highlightRegex, "$1");

  if (searchTerm === "") {
    return;
  }

  // 递归遍历所有文本节点并高亮搜索词
  function highlightText(node) {
    if (node.nodeType === 3) {
      // 文本节点
      const text = node.textContent;
      const regex = new RegExp(searchTerm, "gi");
      if (regex.test(text)) {
        const span = document.createElement("span");
        span.innerHTML = text.replace(
          regex,
          (match) => `<span class="search-highlight">${match}</span>`
        );
        node.parentNode.replaceChild(span, node);
      }
    } else if (node.nodeType === 1) {
      // 元素节点
      // 跳过已经高亮的元素
      if (!node.classList?.contains("search-highlight")) {
        Array.from(node.childNodes).forEach((child) => highlightText(child));
      }
    }
  }

  // 保存原有样式
  const originalStyles = {};
  container.querySelectorAll("*").forEach((el) => {
    originalStyles[el] = el.getAttribute("style");
  });

  // 添加新的高亮
  highlightText(container);

  // 恢复原有样式
  container.querySelectorAll("*").forEach((el) => {
    if (originalStyles[el]) {
      el.setAttribute("style", originalStyles[el]);
    }
  });
}
