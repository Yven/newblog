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

function resetSubmitBtn(type) {
  let form = document.getElementById(type + "Form");
  const button = form.querySelector('button[type="submit"]');
  button.disabled = false;
  button.style.cursor = "pointer";
  button.classList.remove("bg-blue-900/50");
  button.classList.add("bg-blue-900");
  button.classList.add("hover:bg-primary/90");
  const loadingIcon = document.getElementById(type + "Loading");
  loadingIcon.classList.add("hidden");
}

function openModal(type, content, submitFunc) {
  const modal = document.getElementById(type + "Modal");
  if (type == "std" && content != undefined && content !== "") {
    document.getElementById(type + "Content").innerHTML = content;
  }

  let form = document.getElementById(type + "Form");
  // 清除之前可能存在的submit事件监听器
  const oldSubmitListeners = form.cloneNode(true);
  form.parentNode.replaceChild(oldSubmitListeners, form);
  form = document.getElementById(type + "Form");

  form.addEventListener("submit", (e) => {
    const button = form.querySelector('button[type="submit"]');
    button.disabled = true;
    button.style.cursor = "not-allowed";
    button.classList.remove("bg-blue-900");
    button.classList.add("bg-blue-900/50");
    button.classList.remove("hover:bg-primary/90");

    const loadingIcon = document.getElementById(type + "Loading");
    loadingIcon.classList.remove("hidden");

    submitFunc().then(() => {
      resetSubmitBtn(type);
    });

    e.preventDefault();
  });

  const cancelModal = document.getElementById(type + "Cancel");
  cancelModal.addEventListener("click", () => {
    closedModal(type);
  });

  modal.classList.remove("hidden");
}

function closedModal(eleId) {
  return new Promise((resolve) => {
    const modal = document.getElementById(eleId + "Modal");
    modal.classList.add("hidden");

    let form = document.getElementById(eleId + "Form");
    form.reset();

    resetSubmitBtn(eleId);

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
      ? "nav-link text-nowrap sm:text-base text-sm text-gray-600 hover:text-primary dark:text-gray-400 dark:hover:text-gray-100"
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

function getRoute() {
  return window.location.hash.slice(1);
}
