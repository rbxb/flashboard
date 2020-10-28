const maxPostSize = 200;
const postInput = document.querySelector("#post-input");
const postWarning = document.querySelector("#post-warning");
const postButton = document.querySelector("#post-button");
const postList = document.querySelector("#posts-list");
const postTemplate = document.querySelector("#post-template");

function clearPosts() {
	postList.innerHTML = "";
}

function getPosts() {
	clearPosts();
	const req = new XMLHttpRequest();
	req.onreadystatechange = function() {
		if (this.readyState == 4 && this.status == 200) {
			const posts = JSON.parse(this.responseText);
			for (var i = 0; i < posts.length; i++) {
				const clone = postTemplate.content.cloneNode(true);
				const text = window.atob(posts[posts.length - i - 1]);
				clone.querySelector("input").value = text;
				postList.appendChild(clone);
			}
		}
	};
	req.open("GET", "api", true);
	req.send();
}

function sendPost() {
	if (postButton.getAttribute("ok") == "true") {
		const text = postInput.value;
		const req = new XMLHttpRequest();
		req.onreadystatechange = function() {
			if (this.readyState == 4) {
				getPosts();
				if (this.status == 200) {
					postInput.value = "";
					inputChange();
				} else {
					alert("Couldn't post to the server.");
				}
			}
		};
		req.open("POST", "api", true);
		req.send(text);
		postInput.blur();
		postButton.setAttribute("ok", false);
	}
}

function copyText(el) {
	el.parentElement.querySelector("input").select();
	document.execCommand("copy");
	if (window.getSelection) {
		window.getSelection().removeAllRanges();
	} else if (document.selection) {
		document.selection.empty();
	}
	el.parentElement.style.animation = "copyHighlight 0.4s 1";
	setTimeout(function(){el.parentElement.style.animation = "";}, 400);
}

function inputChange() {
	const len = postInput.value.length
	postWarning.innerHTML = len + " / " + maxPostSize;
	if (len <= maxPostSize && len > 0) {
		postButton.setAttribute("ok", true);
	} else {
		postButton.setAttribute("ok", false);
	}
}

inputChange();
getPosts();