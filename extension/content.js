//Websocket
var websocketArguments = 'wss://localhost:8080/sh';
var webSocket;

function onError(error)
{
   console.log(`Error: ${error}`);
}

function createWebsocket()
{
   webSocket = new WebSocket(websocketArguments);
   webSocket.onerror = onWebSocketError;
   webSocket.onopen = onWebSocketOpen;
   webSocket.onmessage = onWebSocketMessage;
}

function onWebSocketError(event)
{
   console.log("WebSocket error observed:", event);
};

function onWebSocketOpen(event)
{
   console.log("console.sh WebSocket open");
};

function onWebSocketMessage(event)
{
   $("#console-extension #console-result").html(event.data);
};

// Workaround to capture Esc key on certain sites
var isOpen = false;
document.onkeyup = (e) => {
	if (e.key == "Escape" && isOpen) {
		browser.runtime.sendMessage({request:"close-console"})
	}
}

$(document).ready(() => {
	createWebsocket()
	var actions = [];
	var isFiltered = false;

	// Append the console into the current page
	$.get(browser.runtime.getURL('/content.html'), (data) => {
		$(data).appendTo('body');


		// New tab page workaround
		if (window.location.href == "browser-extension://mpanekjjajcabgnlbabmopeenljeoggm/newtab.html") {
			isOpen = true;
			$("#console-extension").removeClass("console-closing");
			window.setTimeout(() => {
				$("#console-extension input").focus();
			}, 100);
		}
	});


	// Add actions to the console
	function populateconsole() {
		$("#console-extension #console-result").html("");

	}

	// Open the console
	function openconsole() {
		browser.runtime.sendMessage({request:"get-actions"}, (response) => {
			isOpen = true;
			actions = response.actions;
			$("#console-extension input").val("");
			populateconsole();
			$("html, body").stop();
			$("#console-extension").removeClass("console-closing");
			window.setTimeout(() => {
				$("#console-extension input").focus();
				focusLock.on($("#console-extension input").get(0));
				$("#console-extension input").focus();
			}, 100);
		});
	}

	// Close the console
	function closeconsole() {
		if (window.location.href == "browser-extension://mpanekjjajcabgnlbabmopeenljeoggm/newtab.html") {
			browser.runtime.sendMessage({request:"restore-new-tab"});
		} else {
			isOpen = false;
			$("#console-extension").addClass("console-closing");
		}
	}

	// Hover over an action in the console
	function hoverItem() {
		$(".console-item-active").removeClass("console-item-active");
		$(this).addClass("console-item-active");
	}


	// Autocomplete commands. Since they all start with different letters, it can be the default behavior
	function checkShortHand(e, value) {
		var el = $(".console-extension input");
		if (e.keyCode != 8) {
			if (value == "/t") {
				el.val("/tabs ")
			} else if (value == "/b") {
				el.val("/bookmarks ")
			} else if (value == "/h") {
				el.val("/history ");
			} else if (value == "/r") {
				el.val("/remove ");
			} else if (value == "/a") {
				el.val("/actions ");
			}
		} else {
			if (value == "/tabs" || value == "/bookmarks" || value == "/actions" || value == "/remove" || value == "/history") {
				el.val("");
			}
		}
	}

	// Add protocol
	function addhttp(url) {
			if (!/^(?:f|ht)tps?\:\/\//.test(url)) {
					url = "http://" + url;
			}
			return url;
	}


	// Search for an action in the console
	function search(e) {
		if (e.keyCode == 37 || e.keyCode == 38 || e.keyCode == 39 || e.keyCode == 40 || e.keyCode == 13 || e.keyCode == 37) {
			return;
		}
		var value = $(this).val().toLowerCase();
		checkShortHand(e, value);
		value = $(this).val().toLowerCase();
		if (value.startsWith("/history")) {
			$(".console-item[data-index='"+actions.findIndex(x => x.action == "search")+"']").hide();
			$(".console-item[data-index='"+actions.findIndex(x => x.action == "goto")+"']").hide();
			var tempvalue = value.replace("/history ", "");
			var query = "";
			if (tempvalue != "/history") {
				query = value.replace("/history ", "");
			}
			browser.runtime.sendMessage({request:"search-history", query:query}, (response) => {
				populateconsoleFilter(response.history);
			});
		} else if (value.startsWith("/bookmarks")) {
			$(".console-item[data-index='"+actions.findIndex(x => x.action == "search")+"']").hide();
			$(".console-item[data-index='"+actions.findIndex(x => x.action == "goto")+"']").hide();
			var tempvalue = value.replace("/bookmarks ", "");
			if (tempvalue != "/bookmarks" && tempvalue != "") {
				var query = value.replace("/bookmarks ", "");
				browser.runtime.sendMessage({request:"search-bookmarks", query:query}, (response) => {
					populateconsoleFilter(response.bookmarks);
				});
			} else {
				populateconsoleFilter(actions.filter(x => x.type == "bookmark"));
			}
		} else {
			if (isFiltered) {
				populateconsole();
				isFiltered = false;
			}
			$(".console-extension #console-result .console-item").filter(function(){
				if (value.startsWith("/tabs")) {
					$(".console-item[data-index='"+actions.findIndex(x => x.action == "search")+"']").hide();
					$(".console-item[data-index='"+actions.findIndex(x => x.action == "goto")+"']").hide();
					var tempvalue = value.replace("/tabs ", "");
					if (tempvalue == "/tabs") {
						$(this).toggle($(this).attr("data-type") == "tab");
					} else {
						tempvalue = value.replace("/tabs ", "");
						$(this).toggle(($(this).find(".console-item-name").text().toLowerCase().indexOf(tempvalue) > -1 || $(this).find(".console-item-desc").text().toLowerCase().indexOf(tempvalue) > -1) && $(this).attr("data-type") == "tab");
					}
				} else if (value.startsWith("/remove")) {
					$(".console-item[data-index='"+actions.findIndex(x => x.action == "search")+"']").hide();
					$(".console-item[data-index='"+actions.findIndex(x => x.action == "goto")+"']").hide();
					var tempvalue = value.replace("/remove ", "")
					if (tempvalue == "/remove") {
						$(this).toggle($(this).attr("data-type") == "bookmark" || $(this).attr("data-type") == "tab");
					} else {
						tempvalue = value.replace("/remove ", "");
						$(this).toggle(($(this).find(".console-item-name").text().toLowerCase().indexOf(tempvalue) > -1 || $(this).find(".console-item-desc").text().toLowerCase().indexOf(tempvalue) > -1) && ($(this).attr("data-type") == "bookmark" || $(this).attr("data-type") == "tab"));
					}
				} else if (value.startsWith("/actions")) {
					$(".console-item[data-index='"+actions.findIndex(x => x.action == "search")+"']").hide();
					$(".console-item[data-index='"+actions.findIndex(x => x.action == "goto")+"']").hide();
					var tempvalue = value.replace("/actions ", "")
					if (tempvalue == "/actions") {
						$(this).toggle($(this).attr("data-type") == "action");
					} else {
						tempvalue = value.replace("/actions ", "");
						$(this).toggle(($(this).find(".console-item-name").text().toLowerCase().indexOf(tempvalue) > -1 || $(this).find(".console-item-desc").text().toLowerCase().indexOf(tempvalue) > -1) && $(this).attr("data-type") == "action");
					}
				} else {
					$(this).toggle($(this).find(".console-item-name").text().toLowerCase().indexOf(value) > -1 || $(this).find(".console-item-desc").text().toLowerCase().indexOf(value) > -1);
					if (value == "") {
						$(".console-item[data-index='"+actions.findIndex(x => x.action == "search")+"']").hide();
						$(".console-item[data-index='"+actions.findIndex(x => x.action == "goto")+"']").hide();
					} else {
						$(".console-item[data-index='"+actions.findIndex(x => x.action == "search")+"']").hide();
						$(".console-item[data-index='"+actions.findIndex(x => x.action == "goto")+"']").show();
						$(".console-item[data-index='"+actions.findIndex(x => x.action == "goto")+"'] .console-item-name").html(value);
					}
				}
			});
		}
		
		$(".console-item-active").removeClass("console-item-active");
		$(".console-extension #console-result .console-item:visible").first().addClass("console-item-active");
	}

	// Handle actions from the console
	function handleAction(e) {
		var action = actions[$(".console-item-active").attr("data-index")];
		closeconsole();
		if ($(".console-extension input").val().toLowerCase().startsWith("/remove")) {
			browser.runtime.sendMessage({request:"remove", type:action.type, action:action});
		} else if ($(".console-extension input").val().toLowerCase().startsWith("/history")) {
			if (e.ctrlKey || e.metaKey) {
				window.open($(".console-item-active").attr("data-url"));
			} else {
				window.open($(".console-item-active").attr("data-url"), "_self");
			}
		} else if ($(".console-extension input").val().toLowerCase().startsWith("/bookmarks")) {
			if (e.ctrlKey || e.metaKey) {
				window.open($(".console-item-active").attr("data-url"));
			} else {
				window.open($(".console-item-active").attr("data-url"), "_self");
			}
		} else {
			browser.runtime.sendMessage({request:action.action, tab:action, query:$(".console-extension input").val()});
			switch (action.action) {
				case "bookmark":
					if (e.ctrlKey || e.metaKey) {
						window.open(action.url);
					} else {
						window.open(action.url, "_self");
					}
					break;
				case "scroll-bottom":
					window.scrollTo(0,document.body.scrollHeight);
					break;
				case "scroll-top":
					window.scrollTo(0,0);
					break;
				case "navigation":
					if (e.ctrlKey || e.metaKey) {
						window.open(action.url);
					} else {
						window.open(action.url, "_self");
					}
					break;
				case "fullscreen":
					var elem = document.documentElement;
					elem.requestFullscreen();
					break;
				case "new-tab":
					window.open("");
					break;
				case "email":
					window.open("mailto:");
					break;
				case "url":
					if (e.ctrlKey || e.metaKey) {
						window.open(action.url);
					} else {
						window.open(action.url, "_self");
					}
					break;
				case "goto":
					if (e.ctrlKey || e.metaKey) {
						window.open(addhttp($(".console-extension input").val()));
					} else {
						window.open(addhttp($(".console-extension input").val()), "_self");
					}
					break;
				case "print":
					window.print();
					break;
				case "remove-all":
				case "remove-history":
				case "remove-cookies":
				case "remove-cache":
				case "remove-local-storage":
				case "remove-passwords":
					break;
			}
		}

		// Fetch actions again
		browser.runtime.sendMessage({request:"get-actions"}, (response) => {
			actions = response.actions;
			populateconsole();
		});
	}

	// Customize the shortcut to open the console box
	function openShortcuts() {
		browser.runtime.sendMessage({request:"extensions/shortcuts"});
	}


	// Check which keys are down
	var down = [];

	$(document).keydown((e) => {
		down[e.keyCode] = true;
		if (down[38]) {
			// Up key
			if ($(".console-item-active").prevAll("div").not(":hidden").first().length) {
				var previous = $(".console-item-active").prevAll("div").not(":hidden").first();
				$(".console-item-active").removeClass("console-item-active");
				previous.addClass("console-item-active");
				previous[0].scrollIntoView({block:"nearest", inline:"nearest"});
			}
		} else if (down[40]) {
			// Down key
			if ($(".console-item-active").nextAll("div").not(":hidden").first().length) {
				var next = $(".console-item-active").nextAll("div").not(":hidden").first();
				$(".console-item-active").removeClass("console-item-active");
				next.addClass("console-item-active");
				next[0].scrollIntoView({block:"nearest", inline:"nearest"});
			}
		} else if (down[27] && isOpen) {
			// Esc key
			closeconsole();
		} else if (down[13] && isOpen) {
			// Enter key
			command = $(".console-extension input").val()
			webSocket.send(command)
			//handleAction(e);
		}
	}).keyup((e) => {
		if (down[18] && down[16] && down[80]) {
			if (actions.find(x => x.action == "pin") != undefined) {
				browser.runtime.sendMessage({request:"pin-tab"});
			} else {
				browser.runtime.sendMessage({request:"unpin-tab"});
			}
			browser.runtime.sendMessage({request:"get-actions"}, (response) => {
				actions = response.actions;
				populateconsole();
			});
		} else if (down[18] && down[16] && down[77]) {
			if (actions.find(x => x.action == "mute") != undefined) {
				browser.runtime.sendMessage({request:"mute-tab"});
			} else {
				browser.runtime.sendMessage({request:"unmute-tab"});
			}
			browser.runtime.sendMessage({request:"get-actions"}, (response) => {
				actions = response.actions;
				populateconsole();
			});
		} else if (down[18] && down[16] && down[67]) {
			window.open("mailto:");
		}

		down = [];
	});

	// Recieve messages from background
	browser.runtime.onMessage.addListener((message, sender, sendResponse) => {
		if (message.request == "open-console") {
			if (isOpen) {
				closeconsole();
			} else {
				openconsole();
			}
		} else if (message.request == "close-console") {
			closeconsole();
		}
	});

	$(document).on("click", "#open-page-console-extension-thing", openShortcuts);
	$(document).on("mouseover", ".console-extension .console-item:not(.console-item-active)", hoverItem);
	$(document).on("keyup", ".console-extension input", search);
	$(document).on("click", ".console-item-active", handleAction);
});
