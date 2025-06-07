// marked设置
const { gfmHeadingId } = globalThis.markedGfmHeadingId;
const { mangle } = globalThis.markedMangle;
marked.use(gfmHeadingId({ prefix: "yven-header-" }), mangle(), {
  gfm: true,
});

const indexTpl = `
<div class="py-2">
  <h3 class="sm:text-lg text-md font-semibold text-primary mb-3">
    {{year}}
    {{item}}
  </h3>
</div>
`;

const indexItemTpl = `
<div class="ml-4 py-2 space-y-6">
  <div class="flex justify-between items-center">
    <span class="text-nowrap text-sm text-gray-400 dark:text-gray-600 mr-3">{{category}}</span>
    <a href="#{{slug}}" class="text-gray-800 dark:text-gray-400 hover:text-primary dark:hover:text-gray-300 flex items-center">
      {{title}}
    </a>
    <div class="flex-1 border-b border-dashed border-gray-400 mx-2"></div>
    <span class="text-nowrap text-sm text-gray-400 dark:text-gray-400">{{date}}</span>
  </div>
</div>
`;

const navLinkTpl = `
<a
  href="{{path}}"
  target="{{blank}}"
  class="nav-link text-gray-600 hover:text-primary dark:text-gray-300 dark:hover:text-gray-100"
  >{{title}}</a
>
`;
const navDelimiterTpl = `
<span class="nav-delimiter text-gray-600 dark:text-gray-300">|</span>
`;

const loadingTpl = `
<div class="relative overflow-hidden bg-gray-300 dark:bg-gray-600 rounded h-{{height}} w-{{width}}">
  <div class="absolute inset-0 -translate-x-full animate-[shimmer_1.5s_infinite] bg-gradient-to-r from-transparent via-white/60 to-transparent"></div>
</div>
`;

var originalContent;
var webOpen = false;

function pageClaen() {
  originalContent = "";
  showEdit(false);
}

// 根据路由激活页面内容
function render() {
  if (webOpen === false) {
    showDefault();
    return;
  }

  pageClaen();

  let hash = getRoute();
  switch (hash) {
    case "":
    case "index":
    case "home":
      titleAnimate(true);
      setupList();
      break;
    default:
      titleAnimate(false);
      setupContent(hash);
      break;
  }
}

async function baseInfo() {
  return getWebInfo()
    .then((data) => {
      if (data.code === 200) {
        webOpen = data.data.open;

        const webTitle = document.getElementById("webTitle");
        webTitle.innerHTML = data.data.title;
        document.title = data.data.title;
        const webDesc = document.getElementById("webDesc");
        webDesc.innerHTML = data.data.desc;

        let navLinkEle = [];
        let navList = data.data.nav_list;
        // 判断链接是否为外部链接,设置target属性
        navList = navList.map((item) => {
          if (
            item.path.startsWith("http://") ||
            item.path.startsWith("https://")
          ) {
            item.blank = "blank";
          } else {
            item.blank = "";
          }
          return item;
        });
        navList.forEach((item) => {
          navLinkEle.push(buildTpl(navLinkTpl, item));
        });

        const webNav = document.getElementById("webNav");
        webNav.innerHTML = navLinkEle.join(navDelimiterTpl);
      } else {
        throw new Error(data.message);
      }
    })
    .catch((error) => {
      throw error;
    });
}

function init() {
  // 显示标题加载
  showTitleLoading();
  // 显示正文加载
  showLoadding();

  baseInfo()
    .then(() => {
      // 渲染页面
      render();
      // 设置编辑按钮
      setupEdit();
      // 设置删除按钮
      setupDel();
      // 设置恢复按钮
      setupRecover();
      // 设置完全删除按钮
      setupRealDel();
    })
    .catch((error) => {
      document.getElementById("webTitle").innerHTML = "加载失败";
      document.getElementById("webDesc").innerHTML = "管理员可能提桶跑路了";
      document.getElementById("webNav").innerHTML = "";

      showMsg("加载失败," + error.message);
      console.error(error);
      showDefault();
    });
}

function search(searchTerm) {
  if (isHome()) {
    setupList(searchTerm).then(() => {
      highlightContent("home", searchTerm);
    });
  } else {
    highlightContent("markdownContent", searchTerm);
  }
}

function setupModal(type) {
  const modal = document.getElementById(type + "Modal");
  const closeModal = document.getElementById(type + "Close");
  closeModal.addEventListener("click", () => {
    closedModal(type);
  });
  modal.addEventListener("click", (e) => {
    if (e.target === modal) {
      closedModal(type);
    }
  });
  const cancelModal = document.getElementById(type + "Cancel");
  cancelModal.addEventListener("click", () => {
    closedModal(type);
  });
}

async function setupList(keyword) {
  return getList(keyword)
    .then((res) => {
      let data = res.data;

      if (!data || data.length === 0) {
        showDefault();
        return;
      }

      let html = "";
      data.forEach((item) => {
        let subhtml = "";

        item.item.forEach((item) => {
          item.title +=
            item.delete_time === null
              ? ""
              : "<span class='ml-1 text-xs text-red-800'>[已删除]</span>";
          subhtml += buildTpl(indexItemTpl, item);
        });

        item.item = subhtml;

        html += buildTpl(indexTpl, item);
      });

      document.getElementById("home").innerHTML = html;
      showHome();
    })
    .catch((error) => {
      showMsg("加载失败");
      console.error(error);
      showDefault();
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

function showEdit(show) {
  const markdownContent = document.getElementById("markdownContent");
  const editTextarea = document.getElementById("editTextarea");
  const editContent = document.getElementById("editContent");

  if (show) {
    // 禁用页面滚动
    document.body.style.overflow = "hidden";
    editContent.classList.remove("hidden");
    markdownContent.classList.add("hidden");
    markdownContent.innerHTML = "";
    editTextarea.innerHTML = originalContent;
  } else {
    // 恢复页面滚动
    document.body.style.overflow = "auto";
    editContent.classList.add("hidden");
    editTextarea.innerHTML = "";
    markdownContent.classList.remove("hidden");
    renderContent(originalContent);
  }
}

function setupEdit() {
  const editButton = document.getElementById("editButton");
  if (!isLogin()) {
    editButton.classList.add("hidden");
    return;
  }

  const editForm = document.getElementById("editForm");
  const cancelEdit = document.getElementById("cancelEdit");
  const editContent = document.getElementById("editContent");

  editButton.addEventListener("click", () => {
    showEdit(editContent.classList.contains("hidden"));
  });
  cancelEdit.addEventListener("click", () => {
    showEdit(false);
  });
  editForm.addEventListener("submit", (e) => {
    e.preventDefault();
    const formData = new FormData(editForm);
    editedContent(getRoute(), formData).then((data) => {
      if (data.code === 200) {
        refresh();
      } else {
        showMsg(data.message);
      }
    });
  });
}

function setupDel() {
  const btn = document.getElementById("deleteButton");
  if (!isLogin()) {
    btn.classList.add("hidden");
    return;
  }

  btn.addEventListener("click", () => {
    openModal("std", "是否确认删除此文章？", async function () {
      deleteContent(getRoute()).then((data) => {
        if (data.code === 200) {
          closedModal("std").then(() => {
            window.location.href = "#";
          });
        } else {
          showMsg(data.message);
          closedModal("std");
        }
      });
    });
  });
}

function setupRecover() {
  const btn = document.getElementById("recoverButton");
  if (!isLogin()) {
    btn.classList.add("hidden");
    return;
  }

  btn.addEventListener("click", () => {
    openModal("std", "是否确认恢复此文章？", async function () {
      recoverContent(getRoute()).then((data) => {
        if (data.code === 200) {
          closedModal("std").then(() => {
            window.location.href = "#";
          });
        } else {
          showMsg(data.message);
          closedModal("std");
        }
      });
    });
  });
}

function setupRealDel() {
  const btn = document.getElementById("realDeleteButton");
  if (!isLogin()) {
    btn.classList.add("hidden");
    return;
  }

  btn.addEventListener("click", () => {
    openModal(
      "std",
      "是否确认删除此文章？<br/><p class='text-red-700'>此操作不可恢复！</p>",
      async function () {
        realDeleteContent(getRoute()).then((data) => {
          if (data.code === 200) {
            closedModal("std").then(() => {
              window.location.href = "#";
            });
          } else {
            showMsg(data.message);
            closedModal("std");
          }
        });
      }
    );
  });
}

function renderContent(content) {
  const markdownContent = document.getElementById("markdownContent");
  let htmlContent = marked.parse(content);
  // 渲染数学公式
  htmlContent = renderMath(htmlContent);
  const regex = new RegExp(/<a /g, "gi");
  htmlContent = htmlContent.replace(
    regex,
    `<a target="_blank" rel="noreferrer noopener nofollow" `
  );
  markdownContent.innerHTML = htmlContent;

  // 设置复制按钮
  setupCodeCopyButtons();
  // 设置正文目录
  setupToc();
  // 代码段高亮
  hljs.highlightAll();
}

function setupContent(route) {
  getContent(route)
    .then((data) => {
      if (data.code !== 200) {
        showMsg(data.message);
        showDefault();
      } else {
        const title = document.getElementById("title");
        title.innerHTML = data.data.title;
        document.title = data.data.title;
        const time = document.getElementById("time");
        time.innerHTML = data.data.create_time;

        originalContent = data.data.content;
        // 显示正文
        renderContent(originalContent);

        showContent();

        if (isLogin()) {
          // 是否显示按钮
          if (data.data.delete_time == null) {
            document.getElementById("recoverButton").classList.add("hidden");
            document.getElementById("deleteButton").classList.remove("hidden");
            document.getElementById("realDeleteButton").classList.add("hidden");
          } else {
            document.getElementById("recoverButton").classList.remove("hidden");
            document.getElementById("deleteButton").classList.add("hidden");
            document
              .getElementById("realDeleteButton")
              .classList.remove("hidden");
          }
        }
      }
    })
    .catch((error) => {
      showMsg("加载失败");
      console.error(error);
      showDefault();
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

  loginButton.addEventListener("click", () => {
    openModal("login", "", async function () {
      const loginForm = document.getElementById("loginForm");
      const formData = new FormData(loginForm);
      return login(formData)
        .then((data) => {
          if (data.code !== 200 || data.data.token === undefined) {
            showMsg(data.message);
          } else {
            setCookie("token", data.data.token, data.data.exp);
            loginButton.classList.add("hidden");
            logoutButton.classList.remove("hidden");

            closedModal("login").then(() => {
              refresh();
            });
          }
        })
        .catch((error) => {
          showMsg("登录失败");
          console.error(error);
        });
    });
  });

  logoutButton.addEventListener("click", () => {
    openModal("std", "是否确认退出登录？", async function () {
      logout()
        .then((data) => {
          if (data.code === 200) {
            loginButton.classList.remove("hidden");
            logoutButton.classList.add("hidden");

            closedModal("std").then(() => {
              refresh();
            });
          } else {
            showMsg(data.message);
            closedModal("std");
          }
        })
        .catch((error) => {
          showMsg("请求失败");
          closedModal("std");
          console.error(error);
        });

      deleteCookie("token");
    });
  });
}

function renderMath(htmlContent) {
  // 渲染数学公式
  const latexRegex = /\$\$([\s\S]*?)\$\$|\$((?!\$).*?)\$/g;
  const matches = htmlContent.match(latexRegex);

  if (matches) {
    matches.forEach((match) => {
      const isBlock = match.startsWith("$$");
      const formula = isBlock ? match.slice(2, -2) : match.slice(1, -1);
      const encodedFormula = encodeURIComponent(formula);
      const imgUrl = `https://latex.codecogs.com/svg.image?${encodedFormula}`;

      const imgTag = isBlock
        ? `<div class="text-center"><img class="math-img" src="${imgUrl}" alt="${formula}" class="inline-block my-2"/></div>`
        : `<img class="math-img" src="${imgUrl}" alt="${formula}" class="inline-block" style="vertical-align: middle;"/>`;

      htmlContent = htmlContent.replace(match, imgTag);
    });
  }

  return htmlContent;
}

// 监听哈希变化事件
window.addEventListener("hashchange", function () {
  showLoadding();
  cleanSearch();
  render();
});

document.addEventListener("DOMContentLoaded", function () {
  // 设置暗黑模式
  const darkModeToggle = document.getElementById("darkModeToggle");
  const htmlElement = document.documentElement;
  var isDarkMode = window.matchMedia("(prefers-color-scheme: dark)").matches;
  if (localStorage.theme !== undefined) {
    isDarkMode = localStorage.theme === "dark";
  }
  if (isDarkMode) {
    htmlElement.classList.toggle("dark");
    document.getElementById("sunIcon").classList.add("hidden");
    document.getElementById("moonIcon").classList.remove("hidden");
    updateDarkModeStyles(true);
  }
  darkModeToggle.addEventListener("click", function () {
    const isDark = htmlElement.classList.toggle("dark");
    localStorage.theme = isDark ? "dark" : "light";
    if (isDark) {
      document.getElementById("sunIcon").classList.add("hidden");
      document.getElementById("moonIcon").classList.remove("hidden");
    } else {
      document.getElementById("sunIcon").classList.remove("hidden");
      document.getElementById("moonIcon").classList.add("hidden");
    }
    updateDarkModeStyles(isDark);
  });

  // 设置计时时间
  const urodz = new Date("09/06/2022");
  const now = new Date();
  const ile = now.getTime() - urodz.getTime();
  const dni = Math.floor(ile / (1000 * 60 * 60 * 24));
  document.getElementById("timer").innerHTML = dni;

  // 网站信息初始化
  init();

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
  pageFlipContainer.addEventListener("mouseenter", function () {
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
  let clientHeight =
    document.documentElement.clientHeight || document.body.clientHeight;
  window.addEventListener("scroll", function () {
    // 页面中内容的总高度
    let docELHeight = docEl.scrollHeight;
    // 页面内已经滚动的距离
    let scrollTop = docEl.scrollTop;
    // 页面上滚动到底部的条件
    if (scrollTop >= docELHeight - clientHeight) {
      // 页面内已经滚动的距离 = 页面中内容的总高度 - 浏览器可视部分的高度
      pageFlipContainer.dispatchEvent(new MouseEvent("mouseenter"));
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
    if (searchInput.value.trim() === "") {
      searchInput.classList.remove("w-[4.5rem]");
      searchInput.classList.add("w-8");
      searchInput.placeholder = "";
    }
    search(searchInput.value);
  });
  searchInput.addEventListener("keydown", function (e) {
    if (e.key === "Enter") {
      e.preventDefault();
      search(searchInput.value);
      return;
    }
  });

  // 模态框初始化
  setupModal("login");
  setupModal("std");

  // 设置返回顶部按钮
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

  // 设置消息弹窗
  const msgClose = document.getElementById("msgClose");
  const msgModal = document.getElementById("msgModal");
  msgClose.addEventListener("click", () => {
    msgModal.classList.add("hidden");
  });

  // 设置登录登出
  setupLogin();
});
