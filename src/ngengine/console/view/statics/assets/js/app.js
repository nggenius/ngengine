$(function() {
        var $fullText = $('.admin-fullText');
        $('#admin-fullscreen').on('click', function() {
            $.AMUI.fullscreen.toggle();
        });

        $(document).on($.AMUI.fullscreen.raw.fullscreenchange, function() {
            $fullText.text($.AMUI.fullscreen.isFullscreen ? '退出全屏' : '开启全屏');
        });


        var dataType = $('body').attr('data-type');
        for (key in pageData) {
            if (key == dataType) {
                pageData[key]();
            }
        }

        $('.tpl-switch').find('.tpl-switch-btn-view').on('click', function() {
            $(this).prev('.tpl-switch-btn').prop("checked", function() {
                    if ($(this).is(':checked')) {
                        return false
                    } else {
                        return true
                    }
                })
                // console.log('123123123')
        })

        // 请求服务器列表数据
        SendMessage("list", "", InitServerList)
        // 初始化编辑器
        InitEditor()
    })
    // ==========================
    // 侧边导航下拉列表
    // ==========================

$('.tpl-left-nav-link-list').on('click', function() {
        $(this).siblings('.tpl-left-nav-sub-menu').slideToggle(80)
            .end()
            .find('.tpl-left-nav-more-ico').toggleClass('tpl-left-nav-more-ico-rotate');
    })
    // ==========================
    // 头部导航隐藏菜单
    // ==========================

$('.tpl-header-nav-hover-ico').on('click', function() {
    $('.tpl-left-nav').toggle();
    $('.tpl-content-wrapper').toggleClass('tpl-content-wrapper-hover');
})


function InitEditor() {
    var E = window.wangEditor
    var editor = new E('#editor')
    // 或者 var editor = new E( document.getElementById('editor') )
    editor.create()
}

var server_data_array = new Array();
// 初始化服务器列表数据
function InitServerList(result) {
    if (result.status != 200) {
        return
    }
    var info = result.data;
    server_data_array = info.serverdata

    for (var i = 0; i < info.servernum; ++i)
    {
        var tableHtml ="";
        var num = i + 1
        tableHtml += '<tr id="server-data-'+num+'">' +
        '<td>' + info.serverdata[i].name + '</td>' +
        '<td>' + info.serverdata[i].serverip + '</td>' +
        '<td>' + info.serverdata[i].logintime + '</td>' +
        '<td>' + info.serverdata[i].id + '</td>' +
        '<td>' + 
            '<div class="am-btn-group">' +
            '<button type="button" class="am-btn am-btn-primary am-round am-btn-xs" id="btn-start-'+num+'">开启</button>' +
            '<button type="button" class="am-btn am-btn-danger am-round am-btn-xs" id="btn-close-'+num+'">关闭</button>' +
            '<button type="button" class="am-btn am-btn-secondary am-round am-btn-xs" id="btn-restart-'+num+'">重启</button>' +
            '<button type="button" class="am-btn am-btn-warning am-round am-btn-xs" id="btn-config-'+num+'">配置</button>' +
        '</td></tr>';
    }
    var tbody = $(".server-tbody")
    var elements = $(".server-tbody").children().length;  //表示id为“mtTable”的标签下的子标签的个数
    $(".server-tbody").append(tableHtml); //在表头之后添加空白行
    
    BindServerClick(info)
}

// 动态绑定点击事件
function BindServerClick(info) {
    for (var i = 0; i < info.servernum; ++i)
    {
        var num = i + 1
        var servdata = info.serverdata[i]

        var elem = $("#btn-start-" + num)
        elem.bind("click", function() {
            onClickStartBtn(servdata)
        })

        var elem = $("#btn-close-" + num)
        elem.bind("click", function() {
            onClickCloseBtn(servdata)
        })

        var elem = $("#btn-restart-" + num)
        elem.bind("click", function() {
            onClickRestartBtn(servdata)
        })

        var elem = $("#btn-config-" + num)
        elem.bind("click", function() {
            onClickConfigBtn(servdata)
        })
    }
}

function onClickStartBtn(data) {
    alert(data.serverip)
}

function onClickCloseBtn(data) {
    alert(data.serverip)
}

function onClickRestartBtn(data) {
    alert(data.id)
}

function onClickConfigBtn(data) {
    var $dropdown = $('#doc-dropdown-js');
    $dropdown.dropdown('open');
}

function SendMessage(cmd, args, cb) {
    var token = $("input[name='_xsrf']").val();
    var formdata = new FormData();
    formdata.append("cmd", '{"cmd":"' + cmd + '","args":"' + args + '"}');
    formdata.append("_xsrf", token);
    jQuery.ajax({
        type: "post",
        url: "/server_control",
        dataType: "json",
        data: formdata,
        processData: false,     // false序列化data   true不序列化data
        contentType: false,
        /*beforeSend: function (xhr) {
            var token = $("input[name='_xsrf']").val()
            xhr.setRequestHeader('_xsrf', token);
        },*/
        success: cb
    })
}

// 页面数据
var pageData = {
    // ===============================================
    // 首页
    // ===============================================

}