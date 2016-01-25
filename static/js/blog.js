///////////////////////////////////////////

// var code = {
//     Cd_operate_success: "1",  // 操作成功
//     Cd_operate_failed: "-1", // 操作失败
//     Cd_password_error: "5",  // 密码错误
//     Cd_query_failed: "6" // 查询失败
// };
// // 全局变量
// var url = {
//     login: "/login",
// };

// var REG = {
//     viewArt: /^#id[0-9][0-9]?$/
// };

// var PATH = "";

// // －－－－－－－－－－－－－－－－－－－－－－－－－－－－－－－－－－－－－－－－－－－－－－－－－－－－

// $(function () {
//     // 检查是否登录
//     if ($("#login-account").text() == ""){
//         $(".right-bar").css("display", "none")
//         $(".left-bar").css("display", "none");
//     }else{
//         $(".right-bar").css("display", "block")
//         $(".left-bar").css("display", "block");
//     }
//     // 初始化全局变量
//     PATH = location.pathname


//     completeHtml();
//     // document stats
    // document.onreadystatechange = function () {
    //     if (document.readyState == "complete"){
    //         $(".preloader").delay(100).fadeOut("fast")
    //     }
    // };
//     document.onkeydown = function (e) {
//         if (e.keyCode == 13 && $("#password").is(":focus")){
//             siginIn();
//         }
//     };
//     $("#signin").click(function(){
//         siginIn();
//     });
//     $('#logout').off().click(function () {
//         logout();
//     });
//     /* left-bar */
//     $("#open-editor").click(function(event) {
//         $(".editor").fadeIn("fast");
//         $(".mask-pop-layer").fadeIn("fast");
//     });

//     /* right-bar */


//     /* editor */
//     $("#btn-close").click(function(){
//         $(".mask-pop-layer").fadeOut("fast");
//         $(".editor").fadeOut("fast");
//     });

//     $("#btn-send").click(function(){
//         createTopic(false);
//     });
//     $("#btn-save").click(function(){
//         createTopic(true);
//     });
//     $("#btn-top").click(function(){
//         var v = $(this).attr("value")
//         if (v == "true"){
//             $(this).attr("value", "false");
//             $(this).css("color", "black");
//         }else{
//             $(this).attr("value", "true");
//             $(this).css("color", "#1ABC9C");
//         }
//     });
//     // topic
//     getArticle();

//     // 点击事件
//     $('.handle').on('click', 'a', function(e) {
//         var curhash = $(this).attr('href');
//         if ($(this).attr("class") != "title"){
//             location.hash = curhash;
//             parseHash(curhash);
//         }
//     });
// });
// // －－－－－－－－－－－－－－－－－－－－－－－－loading－－－－－－－－－－－－－－－－－－－－－－－－－－－－－
// function completeHtml(){
//     // loading
//     var html_loading = '';
//     html_loading += '<div id="status">';
//     html_loading += '<p class="center-text">Loading the content...';
//     html_loading += '<em>Loading depends on your connection speed!</em>';
//     html_loading += '</p></div>';
//     $(".preloader").html(html_loading);
// }
// // －－－－－－－－－－－－－－－登录－－－－－－－－－－－－－－－－－－－
// function logout(){
//     var resp = get("post", url.login, {flag:"logout"},false);
//     if (resp.Status == code.Cd_operate_success){
//         location.reload();
//     }
// }


// var position = {
//     1:[
//         "51px",
//         "470px"
//     ],
//     2:[
//         "51px",
//         "300px"
//     ]  
// }

// function tips(type, text){
//     $(".help-tip").css("top", position[type][0]);
//     $(".help-tip").css("right", position[type][1]);
//     $(".help-tip").text(text);
//     $(".help-tip").fadeTo("fast",1);
//     setTimeout(function () {
//         $(".help-tip").fadeTo("fast",0);
//     },1000);
// }
// // －－－－－－－－－－－－－－－－－－－－－－－－文章翻页－－－－－－－－－－－－－－－－－－－－－－－－－－－－－
// // home
// var homepage = 1
// var maxpage = 1
// function getArticle(){
//     regindex = /^\/[a-zA-Z]+$/
//     regarticle = /^\/[a-zA-Z]+\/article$/
//     if (regindex.test(PATH)){
//         var html = '';
//         var resp = get("get", location.pathname + "/article?action=default&page="+homepage, null, false)
//         if (resp.Status == code.Cd_operate_success) {
//             html = showpage(resp.Data);
//         }else{
//             alert(resp.Err);
//         }
//         $("#previews").before(html);
//     }else if (regarticle.test(PATH)){
//         reg = /id=[a-zA-Z0-9]+/
//         var r = location.search.match(reg)
//         if (r != null && r.length > 0){
//             var resp = get("get", PATH+"?flag=article&"+ r[0], null, false);

//             if (resp.Status == code.Cd_operate_success){
//                 $("#article-title").text(resp.Data.Prev.Title)
//                 $(".article-content").append(marked(resp.Data.Cont.Content))
//                 $('pre code').each(function(i, block) {
//                     hljs.highlightBlock(block);
//                 });
//                  //在文章中查找title并填充到div AnchorContent中
//                 $(".article-content").find("h2,h3,h4,h5,h6").each(function(i,item){
//                     var tag = $(item).get(0).localName;
//                     $(item).attr("id","wow"+i);
//                     $("#AnchorContent").append('<li><a class="new'+tag+' anchor-link" onclick="return false;" href="#" link="#wow'+i+'">'+(i+1)+" · "+$(this).text()+'</a></li>');
//                     $(".newh2").css("margin-left",0);
//                     $(".newh3").css("margin-left",10);
//                     $(".newh4").css("margin-left",20);
//                     $(".newh5").css("margin-left",30);
//                     $(".newh6").css("margin-left",40);
//                 });
//                 $("#AnchorContentToggle").click(function(){
//                     var text = $(this).html();
//                     if(text=="目录[-]"){
//                     $(this).html("目录[+]");
//                     $(this).attr({"title":"展开"});
//                 }else{
//                     $(this).html("目录[-]");
//                     $(this).attr({"title":"收起"});
//                 }
//                     $("#AnchorContent").toggle();
//                 });
//                 $(".anchor-link").click(function(){
//                     $("html,body").animate({scrollTop: $($(this).attr("link")).offset().top}, 1000);
//                 });
//             }else{
//                 alert(resp.Err)
//             }
//         }
//     }
// }

// function nextpage(page, type){
//     var u = '';
//     switch (type){
//         case "category":{
//             // u = "/article?action=category&page=" + page+"&category="+category;
//             break;
//         }
//         case "time":{

//             break;
//         }
//         case "tag":{

//             break
//         }
//         case "tid":{

//             break;
//         }
//     }
// }

// function prepage(page, type){

// }

// function showpage(data){
//     var html = '';
//     for (var i = 0; i < data.Previews.length; i++){
//         // 文章预览
//         if (i > 2){
//             html += '<hr/>'
//         }
//         html += '<li class="user-update article handle">';
//         html += '<div class="article-content">';
//         html += '<a class="title" target="_back" href="'+ PATH + "/article?flag=get&id=" + data.Previews[i].ID +'">' + data.Previews[i].Title + '📖</a><p>'; // 📖
//         html += marked(data.Previews[i].Content);
//         html += '</div>';
//         html += '<div class="meta">';
//         html += '<a class="reply toggle-comments-display" data-comments-id="comments-14108828" href="#id' + data.Previews[i].ID + '/comments">评论(' + data.Previews[i].Replys.length + ')</a>、';
//         html += '<a href="#group/'+ data.Previews[i].Group +'"><i class="fa fa-book"></i>'+ data.Previews[i].Group +'</a>、';
//         html += '<i class="fa fa-th"></i>';
//         html += '<a style="color:#ccc;" href="#">阅读('+ data.Previews[i].Views +')</a>';
//         html += '<time>11月14日 23:39</time>';
//         html += '</div>';
//         html += '</li>';
//     }
//     return html;
// }

// function parseHash(hash){
    
// }

// editor
// function createTopic(isdrafts){
//     var art = $.getContent();
//     if (art[0] == "" || art[1] == ""){
//         // 请补充完整
//         return
//     }
//     var top = $("#btn-top").attr("value");
//     var group = $("#category").val();
//     var tags = $("#tag").val();
//     var u = PATH + "/article";
//     var resp = get("post", u, {flag:"add",title:art[0],content:art[1],group:group,tags:tags,isdrafts:isdrafts,top:top}, false)
//     if (resp != null){
//         if (resp.Status == code.Cd_operate_success){
//             location.reload()
//         }else{
//             alert(resp.Err)
//         }
//     }
// }
// ------------------------------------- 配置 ----------------------------------------
var success = 1;
var config = {
    INFO: "info",
    WARNING: "warning",
    SUCCESS: "success",
    ALERT: "alert"
}
marked.setOptions({
    renderer: new marked.Renderer(),
    gfm: true,
    tables: true,
    breaks: false,
    pedantic: false,
    sanitize: true,
    smartLists: true,
    smartypants: false,
    //highlight:function(code){
    //    return hljs.highlightAuto(code).value;
    //}
});
// ------------------------------------- 入口 ----------------------------------------
$(function(){
    window.onerror = function(errorMessage, scriptURI, lineNumber,columnNumber,errorObj) {
        alert("出错行号："+lineNumber+"\n错误信息："+errorMessage+"\n出错文件："+scriptURI);
    };
    // preloader
    document.onreadystatechange = function () {
        if (document.readyState == "complete"){
            $(".preloader").delay(100).fadeOut("fast")
        }
    };
    // login
    document.onkeydown = function (e) {
        if (e.keyCode == 13 &&location.pathname=="/login"){
            siginIn();
        }
    };
    $("#login").on("click", function(){
        siginIn();
    });
    // $("#lost").on("click", function(){
    //     pushMessage(config.INFO, "联系方式|请Email至站长邮箱！");
    // });

    // -------------- left-bar -------------
    // 分类
    $('#nav-category').on('click', 'li', function(){
        if (!$(this).hasClass('active') && !$(this).hasClass("tq")) {
        $('#nav-category li').removeClass('active');
        $(this).addClass('active');
        }
    });
    // 搜索
    $('#btn-search').on('click', function(){
        var content = $('#search-content').val();
        console.log(content)
    });
    // -------------- right-content -------------
    // -------------- 文章列表 -------------
    
});

// ------------------------------------- 功能函数 ----------------------------------------
function timer(){
    var now = (new Date()).valueOf();
    if (now - conf.MODIFY_TIME > 600 && now - conf.MODIFY_TIME < 1101){
        $("#highlight-content").html(marked($("#edit-area").val()));
        $('pre code').each(function(i, block) {
            hljs.highlightBlock(block);
        });
    }
}

function siginIn(){
    username = $("#login-name").val();
    password = $("#login-passwd").val();
    if (username == "" || password == ""){
        pushMessage(config.INFO, "参数错误|请补充完整！");
        return
    }
    var resp = get("post", "/login", {username:username,password:password},false);
    if (resp.Status != success){pushMessage(resp.Err.Level, resp.Err.Msg);return;}
    pushMessage(config.SUCCESS, "登录成功|即将跳转到管理面板！")
    window.setInterval(location.assign(resp.Data),1000);
}

// ------------------------------------- 通信 ----------------------------------------
function get(method, url, data, async) 
{
    var resp;
    $.ajax({
        type: method,
        url: url,
        data: data,
        dataType: 'json',
        async: async,
        success: (function(response){
            resp = response;
        })
    });
    return resp;
}
// ------------------------------------- notify ----------------------------------------
function pushMessage(t,mes){
    $.Notify({
        caption: mes.split("|")[0],
        content: mes.split("|")[1],
        type: t
    });
};
var _notify_container = false;
var _notifies = [];

var Notify = {

    _container: null,
    _notify: null,
    _timer: null,

    version: "3.0.0",

    options: {
        icon: '', // to be implemented
        caption: '',
        content: '',
        shadow: true,
        width: 'auto',
        height: 'auto',
        style: false, // {background: '', color: ''}
        position: 'right', //right, left
        timeout: 3000,
        keepOpen: false,
        type: 'default' //default, success, alert, info, warning
    },

    init: function(options) {
        this.options = $.extend({}, this.options, options);
        this._build();
        return this;
    },

    _build: function() {
        var that = this, o = this.options;

        this._container = _notify_container || $("<div/>").addClass("notify-container").appendTo('body');
        _notify_container = this._container;

        if (o.content === '' || o.content === undefined) {return false;}

        this._notify = $("<div/>").addClass("notify");

        if (o.type !== 'default') {
            this._notify.addClass(o.type);
        }

        if (o.shadow) {this._notify.addClass("shadow");}
        if (o.style && o.style.background !== undefined) {this._notify.css("background-color", o.style.background);}
        if (o.style && o.style.color !== undefined) {this._notify.css("color", o.style.color);}

        // add Icon
        if (o.icon !== '') {
            var icon = $(o.icon).addClass('notify-icon').appendTo(this._notify);
        }

        // add title
        if (o.caption !== '' && o.caption !== undefined) {
            $("<div/>").addClass("notify-title").html(o.caption).appendTo(this._notify);
        }
        // add content
        if (o.content !== '' && o.content !== undefined) {
            $("<div/>").addClass("notify-text").html(o.content).appendTo(this._notify);
        }

        // add closer
        var closer = $("<span/>").addClass("notify-closer").appendTo(this._notify);
        closer.on('click', function(){
            that.close(0);
        });

        if (o.width !== 'auto') {this._notify.css('min-width', o.width);}
        if (o.height !== 'auto') {this._notify.css('min-height', o.height);}

        this._notify.hide().appendTo(this._container).fadeIn('slow');
        _notifies.push(this._notify);

        if (!o.keepOpen) {
            this.close(o.timeout);
        }

    },

    close: function(timeout) {
        var self = this;

        if(timeout === undefined) {
            return this._hide();
        }

        setTimeout(function() {
            self._hide();
        }, timeout);

        return this;
    },

    _hide: function() {
        var that = this;

        if(this._notify !== undefined) {
            this._notify.fadeOut('slow', function() {
                $(this).remove();
                _notifies.splice(_notifies.indexOf(that._notify), 1);
            });
            return this;
        } else {
            return false;
        }
    },

    closeAll: function() {
        _notifies.forEach(function(notEntry) {
            notEntry.hide('slow', function() {
                notEntry.remove();
                _notifies.splice(_notifies.indexOf(notEntry), 1);
            });
        });
        return this;
    }
};

$.Notify = function(options) {
    return Object.create(Notify).init(options);
};

$.Notify.show = function(message, title, icon) {
    return $.Notify({
        content: message,
        caption: title,
        icon: icon
    });
};