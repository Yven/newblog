var path = "http://localhost:8980";

async function login(formData) {
    return fetch(path + "/login", {
        method: "POST",
        body: formData,
    }).then((response) => response.json())
}

async function logout() {
  return fetch(path + "/logout", {
    method: "POST",
    headers: {
      Authorization: "Bearer " + getCookie("token"),
    },
  }).then((response) => response.json());
}

async function getList() {
  return fetch(path + "/list").then((response) => response.json());
}

async function getContent(route) {
  return fetch(path + "/content/" + route).then((response) => response.json());
}

async function editedContent(route, formData) {
  return fetch(path + "/content/" + route, {
    method: "POST",
    headers: {
      Authorization: "Bearer " + getCookie("token"),
    },
    body: formData,
  }).then((response) => response.json());
}

async function deleteContent(route) {
  return fetch(path + "/content/" + route, {
    method: "DELETE",
    headers: {
      Authorization: "Bearer " + getCookie("token"),
    },
  }).then((response) => response.json());
}