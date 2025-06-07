// 封装通用的 HTTP 请求方法
async function request(url, options = {}) {
  try {
    const defaultOptions = {
      headers: {
        "Content-Type": "application/json",
        "Authorization": "Bearer " + getCookie("token"),
      },
    };

    if (options.method === "POST" && options.body instanceof FormData) {
      let formData = options.body;
      let formJson = {};
      for (let [key, value] of formData.entries()) {
        formJson[key] = value;
      }
      options.body = JSON.stringify(formJson);
    }

    const response = await fetch(path + url, {
      ...defaultOptions,
      ...options,
    });

    if (!response.ok) {
      throw new Error("网络请求失败");
    }

    return response.json();
  } catch (error) {
    throw error;
  }
}

// 封装 GET 请求
async function get(url, options = {}) {
  return request(url, { method: "GET", ...options });
}

// 封装 POST 请求
async function post(url, data, options = {}) {
  return request(url, {
    method: "POST",
    body: data,
    ...options,
  });
}

// 封装 DELETE 请求
async function del(url, options = {}) {
  return request(url, {
    method: "DELETE",
    ...options,
  });
}

// -----------------

// 获取网站基本信息
async function getWebInfo() {
  return get("/web/info");
}

// 用户登录
async function login(formData) {
  return post("/login", formData);
}

// 用户登出
async function logout() {
  return post("/logout");
}

// 获取文章列表
async function getList(keyword) {
  // 如果有关键词则添加到查询参数中
  const params = keyword ? `?keyword=${encodeURIComponent(keyword)}` : "";
  return get(`/list${params}`);
}

// 获取指定文章内容
async function getContent(route) {
  return get(`/content/${route}`);
}

// 编辑文章内容
async function editedContent(route, formData) {
  return post(`/content/${route}`, formData);
}

// 删除文章
async function deleteContent(route) {
  return del(`/content/${route}`);
}

// 恢复文章
async function recoverContent(route) {
  return get(`/content/recover/${route}`);
}

// 完全删除文章
async function realDeleteContent(route) {
  return del(`/content/delete/${route}`);
}
