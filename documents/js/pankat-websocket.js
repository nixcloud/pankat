$(document).ready(function() {
  function initWebsocket() {
    path = "";
    if (location.protocol == "http:") {
      path += "ws://"
    } else {
      path += "wss://"
    }

    path += location.hostname + ":" + location.port + "/websocket" 

    var ws = new ReconnectingWebSocket(path);

    ws.onopen = function () {
      console.log('Connection opened');
      $('#websocketStatus').removeClass('glyphicon-remove').addClass('glyphicon-ok');

//       setInterval(function() {
//         ws.send("Client writes: Keep alive"  );
//       }, 2000 );
    };

    ws.onmessage = function (msg) {
//       alert("onmessage")
//       var res = JSON.parse(msg.data).updateTicket;
// console.log (msg)
      console.log ("ws.onmessage = ", msg.data)
      console.log ("ws.onmessage = ", "reload")


      if (msg.data === '"reload"') {
        document.location.reload(true);
      }

//       if (!isNaN(res.sequence)) {
//         if (sequence == undefined && res['id'] == undefined) {
//           sequence = res.sequence;
//         } else if (sequence != res.sequence){
//           sequence = res.sequence;
//           ticketCache = {}; 
//           ticketViewCache = {}; //after removing ticket Views are stored in here 
//           autocompleteValues = [];
//           hidden = []; //indices of removed ticketViews
//           filter = {}; //those filters are used to restrict the shown tickets even more
//           
//           Columns = config['lanes'];
//           $('#lanes').empty();
//           initColumns();
//           initView();
//           //initiate the filters
//           $('.filter').keyup();
//     //      location.reload();
//         }
//       }
     
//       if (res['id'] != undefined && !isNaN(res.id)){
//         console.log("newMessage: update Ticket "+ res.id);
//         getAndUpdateTicket(res.id);
//         sequence++;
//       }  
    };
    ws.onclose = function (msg) {
           $('#websocketStatus').removeClass('glyphicon-ok').addClass('glyphicon-remove');

      console.log('Connection closed');
    };
  };
  initWebsocket();
});