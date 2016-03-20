///////////////////////////////////////////
// ------------------------------------- 配置 ----------------------------------------
var success = 1;
var config = {
    INFO: "info",
    WARNING: "warning",
    SUCCESS: "success",
    ALERT: "alert"
}
// marked.setOptions({
//     renderer: new marked.Renderer(),
//     gfm: true,
//     tables: true,
//     breaks: false,
//     pedantic: false,
//     sanitize: true,
//     smartLists: true,
//     smartypants: false,
//     //highlight:function(code){
//     //    return hljs.highlightAuto(code).value;
//     //}
// });
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

// ------------------------------------- 数组操作 ----------------------------------------
Array.prototype.indexOf = function(val) {
    for (var i = 0; i < this.length; i++) {
        if (this[i] == val) return i;
    }
    return -1;
};
Array.prototype.remove = function(val) {
    var index = this.indexOf(val);
    if (index > -1) {
        this.splice(index, 1);
    }
};
// ------------------------------------- json format ----------------------------------------
(function(window) {
    var p = [],
        push = function( m ) { return '\\' + p.push( m ) + '\\'; },
        pop = function( m, i ) { return p[i-1] },
        tabs = function( count ) { return new Array( count + 1 ).join( '\t' ); };

    window.JSONFormat = function( json ) {
        p = [];
        var out = "",
            indent = 0;
        
        // Extract backslashes and strings
        json = json
            .replace( /\\./g, push )
            .replace( /(".*?"|'.*?')/g, push )
            .replace( /\s+/, '' );      
        
        // Indent and insert newlines
        for( var i = 0; i < json.length; i++ ) {
            var c = json.charAt(i);
            
            switch(c) {
                case '{':
                case '[':
                    out += c + "\n" + tabs(++indent);
                    break;
                case '}':
                case ']':
                    out += "\n" + tabs(--indent) + c;
                    break;
                case ',':
                    out += ",\n" + tabs(indent);
                    break;
                case ':':
                    out += ": ";
                    break;
                default:
                    out += c;
                    break;      
            }                   
        }
        
        // Strip whitespace from numeric arrays and put backslashes 
        // and strings back in
        out = out
            .replace( /\[[\d,\s]+?\]/g, function(m){ return m.replace(/\s/g,''); } )
            .replace( /\\(\d+)\\/g, pop ) // strings
            .replace( /\\(\d+)\\/g, pop ); // backslashes in strings
        
        return out;
    };
})(window);
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