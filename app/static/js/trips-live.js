// Live flight updates for the trips page. Subscribes to the server-side SSE
// proxy at /events (which forwards the account's skyvisor-api stream) and
// surfaces disruption messages without a page reload.
(function () {
  var region = document.getElementById("live-updates");
  if (!region || typeof EventSource === "undefined") {
    return;
  }

  var list = region.querySelector("[data-live-list]");
  var status = region.querySelector("[data-live-status]");
  var flightEvents = [
    "flight.cancelled",
    "flight.delayed",
    "gate.changed",
    "flight.updated",
  ];

  function setStatus(text, live) {
    if (!status) return;
    status.textContent = text;
    status.setAttribute("data-state", live ? "live" : "off");
  }

  function announce(event) {
    var payload;
    try {
      payload = JSON.parse(event.data);
    } catch (err) {
      return;
    }
    if (!payload || !payload.message) {
      return;
    }
    region.hidden = false;
    var item = document.createElement("li");
    item.className =
      "flex items-center gap-2 rounded-md bg-background/70 px-3 py-2 text-sm";
    var time = new Date(payload.at || Date.now()).toLocaleTimeString([], {
      hour: "2-digit",
      minute: "2-digit",
    });
    item.innerHTML =
      '<span class="font-mono text-xs text-muted-foreground">' +
      time +
      "</span> <span>" +
      escapeHTML(payload.message) +
      "</span>";
    if (list) {
      list.insertBefore(item, list.firstChild);
      while (list.children.length > 8) {
        list.removeChild(list.lastChild);
      }
    }
  }

  function escapeHTML(value) {
    var div = document.createElement("div");
    div.textContent = value;
    return div.innerHTML;
  }

  var source = new EventSource("/events");
  source.addEventListener("ready", function () {
    region.hidden = false;
    setStatus("Live", true);
  });
  source.onerror = function () {
    setStatus("Reconnecting…", false);
  };
  flightEvents.forEach(function (name) {
    source.addEventListener(name, announce);
  });

  window.addEventListener("beforeunload", function () {
    source.close();
  });
})();
