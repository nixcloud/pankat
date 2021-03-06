var h = virtualDom.h;
var vdiff = virtualDom.diff;
var patch = virtualDom.patch;
var createElement = virtualDom.create;
var websocket;
var timer = 0;
var content_old = "al33vkl3vj3vk3jvj3v3jvj;sjgs ;g  gss  ;sgs ggs ;;";

function init() {
    console.log("init");
    websocket = new WebSocket(`ws://${window.location.hostname}:8080/entry`);
    websocket.onopen = function(evt) { onOpen(evt) };
    websocket.onclose = function(evt) { onClose(evt) };
    websocket.onmessage = function(evt) { onMessage(evt) };
    websocket.onerror = function(evt) { onError(evt) };
}

function onOpen(evt) {
    console.log("WS: connected");
    var message = $('#note').val();
    doSend(message);
}

function onClose(evt) {
    console.log("WS: disconnected");
}

function toModel(tx) {
    return {
        body: tx
    };
}
var html;
var rootNode;

function onMessage(evt) {
    console.log("WS received");
//     var r = jQuery.parseJSON ( ' {"foo": "bar"} ' );
    var r = JSON.parse(evt.data);
    //jQuery("#previewPane").removeClass("flash");



    if (typeof html === 'undefined' ) {
        //console.log('1')
        var str = "<div>" + r.body + "</div>";
        html = eval(dom2hscript.parseHTML(str));
        rootNode = createElement(html);
        $('#previewPane').append(rootNode);


        //var one = 'beep boop',
        //other = 'beep boob blah',
        //color = '',
        //span = null;

        //var diff = JsDiff.diffChars(one, other),
        //    display = document.getElementById('display'),
        //    fragment = document.createDocumentFragment();
        //
        //diff.forEach(function(part){
        //  // green for additions, red for deletions
        //  // grey for common parts
        //  color = part.added ? 'green' :
        //    part.removed ? 'red' : 'grey';
        //  span = document.createElement('span');
        //  span.style.color = color;
        //  span.appendChild(document
        //    .createTextNode(part.value));
        //  fragment.appendChild(span);
        //});
        //display.appendChild(fragment);



    } else {
        //console.log('2')
        var str = "<div>" + r.body + "</div>";
        var h1 = eval(dom2hscript.parseHTML(str));
        var patches = vdiff(html, h1);
        var oldPatch = vdiff(html, h1);

        console.log('old object', oldPatch)

        // FIXME: code below is still buggy and experimental code to modify the patch to add 'yellow' to changed text fields
        //for (const [key, value] of Object.entries(patches)) {
        //  console.log(key)
        //  if (key !== "a" && key !== 'undefined') {
        //    if ("patch" in value) {
        //      if ("text" in value.patch) {
        //        console.log("der text: " + value.patch.text)
        //        //patches[key].patch.text="<b>tut</b>"
        //        // add properteis
        //        if (!("properties" in value.patch)) {
        //          value.patch.properties = {};
        //        }
        //        // add properties.className: "highlight23"
        //        //value.patch.properties.className = "highlight23"
        //        var t = eval(dom2hscript.parseHTML('<span class="highlight23">' + value.patch.text + '</span>'));
        //        //console.log("t:", t)
        //        value.patch = t;
        //      }
        //    }
        //  }
        //}
        //debugger
        //console.log(r.body)
        //rootNode = createElement(html);
        rootNode = patch(rootNode, patches);
    }
}

function onError(evt) {
    console.log(evt.data);
}

function doSend(message) {
    websocket.send(JSON.stringify(toModel(message)));
    //console.log(message);
    if (websocket.readyState !== 1) {
        //
    } else {
        //init()
    }
}

// http://ejohn.org/blog/how-javascript-timers-work/
function resetTimerOnChange() {
    clearTimeout(timer);
    timer = setTimeout(timerExpired, 1500);
}

function timerExpired () {
    var message = $('#note').val();
    if (message !== content_old) {
        content_old = message;
        doSend(message);
    }
}

window.addEventListener("load", init, false);