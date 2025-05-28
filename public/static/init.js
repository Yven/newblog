// marked设置
const { gfmHeadingId } = globalThis.markedGfmHeadingId;
const { mangle } = globalThis.markedMangle;
marked.use(gfmHeadingId({ prefix: "yven-header-" }), mangle(), {
  gfm: true,
});

const indexTpl = `
<div class="py-2">
  <h3 class="text-lg font-semibold text-primary mb-3">
    {{year}}
    {{item}}
  </h3>
</div>
`;

const indexItemTpl = `
<div class="ml-8 py-2 space-y-3">
  <div class="flex justify-between items-center">
    <a href="#{{slug}}" class="text-gray-800 dark:text-gray-300 hover:text-primary dark:hover:text-gray-100">
      <span class="text-sm text-gray-500 dark:text-gray-400 mr-3">{{category}}</span>
      {{title}}
    </a>
    <div class="flex-1 border-b border-dashed border-gray-400 mx-2"></div>
    <span class="text-sm text-gray-500 dark:text-gray-400">{{date}}</span>
  </div>
</div>
`;

var originalContent;
var htmlContent;
var isEdit = false;

// 根据路由激活页面内容
function render() {
  showDefault();
  const hash = window.location.hash.slice(1);
  switch (hash) {
    case "":
    case "index":
    case "home":
      setupList();
      break;
    default:
      setupContent(hash);
      break;
  }
}

function init() {
  // 目录显示按钮初始化
  const tocToggle = document.getElementById("tocToggle");
  tocToggle.addEventListener("click", function () {
    const isVisible = !tocContainer.classList.contains("hidden");
    localStorage.showToc = Number(!isVisible);
    if (isVisible) {
      tocContainer.classList.add("hidden");
    } else {
      tocContainer.classList.remove("hidden");
    }
  });

  // 登录按钮显示动作初始化
  const pageFlipContainer = document.getElementById("pageFlipContainer");
  pageFlipContainer.addEventListener('mouseenter', function() {
    const foldEffect = document.getElementById("foldEffect");
    foldEffect.classList.add("show-flod-effect");
    const loginButton = document.getElementById("loginButton");
    loginButton.classList.add("show-user-login");
    const logoutButton = document.getElementById("logoutButton");
    logoutButton.classList.add("show-user-login");
  });
  pageFlipContainer.addEventListener("mouseleave", function () {
    const foldEffect = document.getElementById("foldEffect");
    foldEffect.classList.remove("show-flod-effect");
    const loginButton = document.getElementById("loginButton");
    loginButton.classList.remove("show-user-login");
    const logoutButton = document.getElementById("logoutButton");
    logoutButton.classList.remove("show-user-login");
  });

  // 滚动到底部自动显示登录按钮
  let docEl = document.documentElement;
  // 浏览器可视部分的高度
  let clientHeight = document.documentElement.clientHeight || document.body.clientHeight;
  window.addEventListener("scroll", function () {
    // 页面中内容的总高度
    let docELHeight = docEl.scrollHeight;
    // 页面内已经滚动的距离
    let scrollTop = docEl.scrollTop;
    // 页面上滚动到底部的条件
    if (scrollTop >= docELHeight - clientHeight) {
      // 页面内已经滚动的距离 = 页面中内容的总高度 - 浏览器可视部分的高度
      pageFlipContainer.dispatchEvent(new MouseEvent('mouseenter'));
    } else {
      pageFlipContainer.dispatchEvent(new MouseEvent("mouseleave"));
    }
  });

  // 搜索按钮初始化
  const searchButton = document.getElementById("searchButton");
  const searchInput = document.getElementById("searchInput");
  searchButton.addEventListener("click", function () {
    if (searchInput.classList.contains("w-8")) {
      searchInput.focus();
      searchInput.classList.remove("w-8");
      searchInput.classList.add("w-[4.5rem]");
      searchInput.placeholder = "搜索";
    }
  });
  searchInput.addEventListener("blur", function () {
    if (searchInput.value === "") {
      searchInput.classList.remove("w-[4.5rem]");
      searchInput.classList.add("w-8");
      searchInput.placeholder = "";
    }
  });
}

function setupList() {
  getList().then((res) => {
    let data = res.data;

    if (!data || data.length === 0) {
      showDefault();
      return;
    }

    let html = "";
    data.forEach((item) => {
      let subhtml = "";

      item.item.forEach((item) => {
        subhtml += indexItemTpl.replace(/{{\s*(\w+)\s*}}/g, (match, key) => {
          return item[key] !== undefined ? item[key] : "";
        });
      });

      html += indexTpl.replace(/{{\s*(\w+)\s*}}/g, (match, key) => {
        if (key === "item") {
          return subhtml;
        } else {
          return item[key] !== undefined ? item[key] : "";
        }
      });
    });

    document.getElementById("home").innerHTML = html;
    showHome();
  });
}

function setupCodeCopyButtons() {
  // 为代码段增加复制按钮
  document.querySelectorAll("pre code").forEach((block) => {
    const wrapper = document.createElement("div");
    wrapper.className = "code-block-wrapper relative";
    const button = document.createElement("button");
    button.className =
      "copy-button bg-gray-700 text-white text-xs px-2 py-1 rounded-sm !rounded-button";
    button.textContent = "复制";
    button.setAttribute("data-code", block.textContent);

    block.parentNode.parentNode.insertBefore(wrapper, block.parentNode);
    wrapper.appendChild(block.parentNode);
    wrapper.appendChild(button);
  });
  document.querySelectorAll(".copy-button").forEach((button) => {
    button.addEventListener("click", function () {
      const code = this.getAttribute("data-code");
      navigator.clipboard.writeText(code).then(() => {
        const originalText = this.textContent;
        this.textContent = "已复制";
        setTimeout(() => {
          this.textContent = originalText;
        }, 2000);
      });
    });
  });
}

function setupToc() {
  // 目录生成
  const headings = document.querySelectorAll(
    ".markdown-body h2, .markdown-body h3"
  );
  const toc = document.getElementById("toc");
  // 清空
  toc.innerHTML = "";
  const tocContainer = document.getElementById("tocContainer");

  if (localStorage.showToc !== undefined) {
    let isVisible = Number(localStorage.showToc);
    if (isVisible) {
      tocContainer.classList.remove("hidden");
    } else {
      tocContainer.classList.add("hidden");
    }
  }

  headings.forEach((heading) => {
    const id = heading.id;
    const text = heading.textContent;
    const level = heading.tagName.toLowerCase();
    const listItem = document.createElement("li");
    listItem.classList.add("toc-item");
    if (level === "h3") {
      listItem.classList.add("pl-4");
    }
    const link = document.createElement("a");
    link.href = `#${id}`;
    link.textContent = text;
    link.classList.add("block", "py-1", "hover:text-primary");
    listItem.appendChild(link);
    toc.appendChild(listItem);
    link.addEventListener("click", function (e) {
      e.preventDefault();
      document.querySelector(this.getAttribute("href")).scrollIntoView({
        behavior: "smooth",
      });
    });
  });
  // 高亮当前阅读位置
  window.addEventListener("scroll", function () {
    const scrollPosition = window.scrollY;
    let currentSection = null;
    headings.forEach((heading) => {
      const sectionTop = heading.offsetTop - 100;
      if (scrollPosition >= sectionTop) {
        currentSection = heading.id;
      }
    });
    document.querySelectorAll(".toc-item").forEach((item) => {
      item.classList.remove("active");
    });
    if (currentSection) {
      const activeItem = document.querySelector(
        `.toc-item a[href="#${currentSection}"]`
      );
      if (activeItem) {
        activeItem.parentElement.classList.add("active");
      }
    }
  });
}

function setupEdit(route) {
  const editButton = document.getElementById("editButton");
  if (!isLogin()) {
    editButton.classList.add("hidden");
    return;
  }

  const editForm = document.getElementById("editForm");
  const cancelEdit = document.getElementById("cancelEdit");
  editButton.addEventListener("click", () => {
    isEdit = !isEdit;
    const markdownContent = document.getElementById("markdownContent");
    const editContent = document.getElementById("editContent");
    const editTextarea = document.getElementById("editTextarea");
    if (isEdit) {
      editContent.classList.remove("hidden");
      markdownContent.classList.add("hidden");
      markdownContent.innerHTML = "";
      editTextarea.innerHTML = originalContent;
    } else {
      editContent.classList.add("hidden");
      editTextarea.innerHTML = "";
      markdownContent.classList.remove("hidden");
      renderContent(originalContent);
    }
  });
  cancelEdit.addEventListener("click", () => {
    isEdit = false;
    editContent.classList.add("hidden");
    editTextarea.innerHTML = "";
    markdownContent.classList.remove("hidden");
    renderContent(originalContent);
  })
  editForm.addEventListener("submit", (e) => {
    e.preventDefault();
    const formData = new FormData(editForm);
    editedContent(route, formData).then((data) => {
      if (data.code === 200) {
        refresh();
      } else {
        showMsg(data.message);
      }
    });
  });
}

function setupDel(route) {
  if (!isLogin()) {
    document.getElementById("deleteButton").classList.add("hidden");
    return;
  }

  initModal("delete", async function () {
    deleteContent(route).then((data) => {
      if (data.code === 200) {
        closedModal("deletModal").then(() => {
          window.location.href = "#";
        });
      } else {
        showMsg(data.message);
        closedModal("deleteModal");
      }
    });
  });
}

function renderContent(content) {
  const markdownContent = document.getElementById("markdownContent");
  htmlContent = marked.parse(content);
  markdownContent.innerHTML = htmlContent

  // 设置复制按钮
  setupCodeCopyButtons();
  // 设置正文目录
  setupToc();
  // 代码段高亮
  hljs.highlightAll();
  // 渲染数学公式
  renderMath();
}

function setupContent(route) {
  getContent(route).then((data) => {
    if (data.code !== 200) {
      showMsg(data.message);
    } else if (data.data == undefined) {
      return;
    } else {
      const title = document.getElementById("title");
      title.innerHTML = data.data.title;
      const time = document.getElementById("time");
      time.innerHTML = data.data.create_time;

      originalContent = data.data.content;
      renderContent(originalContent);

      // 设置编辑按钮
      setupEdit(route);
      // 设置删除按钮
      setupDel(route);

      // 显示正文
      showContent();
    }
  });
}

function setupSearch() {
  const searchInput = document.getElementById("searchInput");
  searchInput.addEventListener("input", function () {
    let searchTerm = this.value.trim();
    const container = document.getElementById("markdownContent");

    // 移除所有高亮
    const highlightRegex = /<span class="search-highlight">(.*?)<\/span>/g;
    const content = container.innerHTML;
    container.innerHTML = content.replace(highlightRegex, '$1');

    if (searchTerm === "") {
      return;
    }

    // 递归遍历所有文本节点并高亮搜索词
    function highlightText(node) {
      if (node.nodeType === 3) { // 文本节点
        const text = node.textContent;
        const regex = new RegExp(searchTerm, "gi");
        if (regex.test(text)) {
          const span = document.createElement("span");
          span.innerHTML = text.replace(
            regex,
            match => `<span class="search-highlight">${match}</span>`
          );
          node.parentNode.replaceChild(span, node);
        }
      } else if (node.nodeType === 1) { // 元素节点
        // 跳过已经高亮的元素
        if (!node.classList?.contains("search-highlight")) {
          Array.from(node.childNodes).forEach(child => highlightText(child));
        }
      }
    }

    // 保存原有样式
    const originalStyles = {};
    container.querySelectorAll("*").forEach(el => {
      originalStyles[el] = el.getAttribute("style");
    });

    // 添加新的高亮
    highlightText(container);

    // 恢复原有样式
    container.querySelectorAll("*").forEach(el => {
      if (originalStyles[el]) {
        el.setAttribute("style", originalStyles[el]);
      }
    });
  });
}

function setupMsgModal() {
  const closeModal3 = document.getElementById("closeModal3");
  const msgModal = document.getElementById("msgModal");
  closeModal3.addEventListener("click", () => {
    msgModal.classList.add("hidden");
  });
}

function setupDarkMode() {
  // 暗黑模式切换
  const darkModeToggle = document.getElementById("darkModeToggle");
  const htmlElement = document.documentElement;
  var isDarkMode = window.matchMedia("(prefers-color-scheme: dark)").matches;
  if (localStorage.theme !== undefined) {
    isDarkMode = localStorage.theme === "dark";
  }
  if (isDarkMode) {
    htmlElement.classList.toggle("dark");
    darkModeToggle.innerHTML = '<i class="ri-moon-line"></i>';
    updateDarkModeStyles(true);
  }
  darkModeToggle.addEventListener("click", function () {
    const isDark = htmlElement.classList.toggle("dark");
    localStorage.theme = isDark ? "dark" : "light";
    darkModeToggle.innerHTML =
      '<i class="ri-' + (isDark ? "moon" : "sun") + '-line"></i>';
    updateDarkModeStyles(isDark);
  });
}

function updateDarkModeStyles(isDark) {
  // 更新正文文本颜色
  const markdownBody = document.querySelector(".markdown-body");
  if (isDark) {
    markdownBody.style.color = "#9ca3af";
  } else {
    markdownBody.style.color = "";
  }
  const content = document.getElementById("markdownContent");
  // 更新链接和文本颜色
  const links = content.querySelectorAll("a");
  links.forEach((link) => {
    if (isDark) {
      link.classList.add("text-gray-300", "hover:text-gray-200");
      link.classList.remove("text-gray-600", "hover:text-primary");
    } else {
      link.classList.remove("text-gray-300", "hover:text-gray-200");
      link.classList.add("text-gray-600", "hover:text-primary");
    }
  });
}

function setupLogin() {
  const loginButton = document.getElementById("loginButton");
  const logoutButton = document.getElementById("logoutButton");
  if (isLogin()) {
    loginButton.classList.add("hidden");
    logoutButton.classList.remove("hidden");
  } else {
    loginButton.classList.remove("hidden");
    logoutButton.classList.add("hidden");
  }

  initModal("login", async function () {
    const formData = new FormData(loginForm);
    return login(formData).then((data) => {
      if (data.code !== 200 || data.data.token === undefined) {
        showMsg(data.message);
      } else {
        setCookie("token", data.data.token, data.data.exp);
        loginButton.classList.add("hidden");
        logoutButton.classList.remove("hidden");

        closedModal("loginModal").then(() => {
          refresh();
        });
      }
    });
  });
  initModal("logout", async function () {
    logout().then((data) => {
      deleteCookie("token");

      if (data.code === 200) {
        loginButton.classList.remove("hidden");
        logoutButton.classList.add("hidden");

        closedModal("logoutModal").then(() => {
          refresh();
        });
      } else {
        showMsg(data.message);
        closedModal("logoutModal");
      }
    });
  });
}

function setupBackToTop() {
  // 返回顶部按钮
  const backToTopButton = document.getElementById("backToTop");
  window.addEventListener("scroll", function () {
    if (window.scrollY > 300) {
      backToTopButton.classList.add("opacity-30");
      backToTopButton.classList.remove("hidden");
    } else {
      backToTopButton.classList.remove("opacity-30");
      backToTopButton.classList.add("hidden");
    }
  });
  backToTopButton.addEventListener("click", function () {
    window.scrollTo({
      top: 0,
      behavior: "smooth",
    });
  });
}

function renderMath() {
  // 渲染数学公式
  renderMathInElement(document.getElementById("markdownContent"), {
    delimiters: [
      { left: "$$", right: "$$", display: true },
      { left: "$", right: "$", display: false },
    ],
    throwOnError: false,
  });
}

// 监听哈希变化事件
window.addEventListener("hashchange", render);

document.addEventListener("DOMContentLoaded", function () {
  init();
  // 渲染页面
  render();
  // 设置搜索框
  setupSearch();
  // 设置暗黑模式
  setupDarkMode();
  // 设置返回顶部按钮
  setupBackToTop();
  // 设置登录登出
  setupLogin();
  // 设置消息弹窗
  setupMsgModal();
});
