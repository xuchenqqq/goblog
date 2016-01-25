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

var conf = {
    EDIT_MODEL: 1,
    MODIFY_TIME: 0
};

$(function () {
    completeEditHtml();
    completeEditCss();
    // timer
    window.setInterval(function () {
        timer();
    }, 500);
    
    $("#edit-area").on("keydown",function (e) { // metaKey->command
        console.log(e.keyCode);
        conf.MODIFY_TIME = (new Date()).valueOf();
        switch (e.keyCode){
            case 9: { // tab
                e.preventDefault();
                var indent = '    ';
                var start = this.selectionStart;
                var end = this.selectionEnd;
                var selected = window.getSelection().toString();
                selected = indent + selected.replace(/\n/g,'\n'+indent);
                this.value = this.value.substring(0,start) + selected + this.value.substring(end);
                this.setSelectionRange(start+indent.length,start+selected.length);
                break;
            }
            case 90: { // ctrl+z
                if (e.metaKey){
                    key_ctrl_z(e);
                    break;
                }
            }
            case 83: { // ctrl+s
                if (e.metaKey){

                    break;
                }
                if (e.metaKey && e.shiftKey){

                }
                break;
            }
        }
    });
    // 定时保存
    saveLog();
    // test
    $("#btn-marked").click(function () {
        if (conf.EDIT_MODEL == 1) {
            conf.EDIT_MODEL = 0;
            $("#btn-marked").removeClass("glyphicon-eye-open");
            $("#btn-marked").addClass("glyphicon-eye-close");
        }else{
            conf.EDIT_MODEL = 1;
            $("#btn-marked").addClass("glyphicon-eye-open");
            $("#btn-marked").removeClass("glyphicon-eye-close");
        }
        changeEditModel();
    });
    // theme
    initStyles();
});


// 全局函数
jQuery.getContent = function(){
    return new Array($("#title").val(),$("#edit-area").val())
}
// end


// scheduler
function timer(){
    var now = (new Date()).valueOf();
    if (now - conf.MODIFY_TIME > 600 && now - conf.MODIFY_TIME < 1101){
        $("#edit-view").html(marked($("#edit-area").val()));
        $('pre code').each(function(i, block) {
            hljs.highlightBlock(block);
        });
    }
}


// key
var log = [];
function key_ctrl_z(e){// ctrl+z
    e.preventDefault();
    log.pop();
    $("#edit-area").val(log[log.length - 1]).blur();
    $("#edit-area").focus();
}

function saveLog(){
    window.setInterval(function () {
        if (log[log.length - 1] != $("#edit-area").val()) {
            log[log.length] = $("#edit-area").val();
        }
    }, 1500);
}

// 完善editor
function completeEditHtml(){
    var html_edit = '';
    html_edit += '<div id="tool-bar">';
    html_edit += '<span id="btn-marked" title="change model" class="glyphicon glyphicon-eye-open" aria-hidden="true"></span>|';// toolbar
    html_edit += '<span class="glyphicon glyphicon-bold" aria-hidden="true"></span>|'; 
    html_edit += '<span class="glyphicon glyphicon-italic" aria-hidden="true"></span>|'; 
    html_edit += '<span class="glyphicon glyphicon-font" aria-hidden="true"></span>|'; 
    html_edit += '<span class="glyphicon glyphicon-align-left" aria-hidden="true"></span>|'; 
    html_edit += '<span class="glyphicon glyphicon-align-center" aria-hidden="true"></span>|'; 
    html_edit += '<span class="glyphicon glyphicon-align-right" aria-hidden="true"></span>|'; 
    html_edit += '<span class="glyphicon glyphicon-link" aria-hidden="true"></span>|';
    html_edit += '<span class="glyphicon glyphicon-flash" aria-hidden="true"></span>|';
    html_edit += '<span class="glyphicon glyphicon-hd-video" aria-hidden="true"></span>|';
    html_edit += '<span class="glyphicon glyphicon-fullscreen" aria-hidden="true"></span>';
    html_edit += '<div id="div-title" class="input-group input-group-sm"><span class="input-group-addon">标题</span><input id="title" type="text" class="form-control" aria-describedby="sizing-addon2"></div>';
    html_edit += '</div>';
    html_edit += '<textarea id="edit-area">';
    html_edit += '</textarea>';
    html_edit += '<div id="edit-view"></div>';
    $("#markdown-editor").html(html_edit);
}

// 转换编辑模式
function changeEditModel(){
    switch (conf.EDIT_MODEL){
        case 1:{ // 即时浏览
            $("#edit-view").css("display","inline-block");
            $("#edit-area").css("width","50%");
            $("#edit-area").css("border-left", "none");
            break;
        }
        default :{ // 普通
            $("#edit-view").css("display","none");
            $("#edit-area").css("width","100%");
            $("#edit-area").css("border-left", "1px solid #ccc")
        }
    }
}

// css
function completeEditCss(){
    // editor
    $("#tool-bar span").css({
        "margin-left": '12px',
        "margin-right": '12px',
    });

    $("#markdown-editor").css({
        "margin": "0",
        "height": "100%",
        "font-family": "'Helvetica Neue', Arial, sans-serif",
        "color": "#333"
    });
    $("#edit-area, #edit-view").css({
        "display": "inline-block",
        "width": "50%",
        "height": "530px",
        "vertical-align": "top",
        "-webkit-box-sizing": "border-box",
        "-moz-box-sizing": "border-box",
        "box-sizing": "border-box",
        "padding": "0 20px",
        "overflow-y": "auto"
    });

    $("#edit-view").css({
        "font-size": '14px',
        "background-color": "white"
    });

    $("#edit-area").css({
        "border": "none",
        "border-right": "1px solid #ccc",
        "resize": "none",
        "outline": "none",
        "background-color": "#f6f6f6",
        "font-family": "'Monaco', courier, monospace",
        "padding": "10px",
        "font-size": "12px"
    });

    $("#tool-bar").css({
        "border-bottom": "1px solid #ccc",
        "padding":"4px"
    });

    $("#tool-bar span").css({
        "cursor":"pointer"
    });

    $("#div-title").css({
        "margin-top": "10px",
    });

    $("#div-title span").css({
        "font-size": "18px",
        // "color": "#1ABC9C"
    });
}
// 更换主题
function selectStyle(style) {
    $('link[title]').each(function(i, link) {
        link.disabled = (link.title != style);
    });
}

function initStyles() {
    var ul = $('#styles');
    $('link[title]').each(function(i, link) {
        ul.append('<li>' + link.title + '</li>');
    });
    $('#styles li').click(function(e) {
        $('#styles li').removeClass('current');
        $(this).addClass('current');
        selectStyle($(this).text());
    });
    $('#styles li:first-child').click();
}
