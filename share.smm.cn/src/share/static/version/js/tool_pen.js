'use strict';
(function(){

    var trade = window.smm.trade;

    var dom = {
        pen      : $('.smm_frame_pen'),
        canvas   : $('.smm_frame_canvas'),
        close    : $('.smm_frame_close'),
        clear    : $('.smm_frame_clear'),
        group    : $('.smm_frame_canvas_group')

    };

    var p = null;
    var status = {
        last_x : 0,
        last_y : 0,
        draw : false
    };

    dom.canvas._init = function(){
        dom.group._resize();
        p = dom.canvas.get(0).getContext('2d');
        p.strokeStyle = "red";
    };
    dom.canvas._mousedown = function (event) {
        status.draw = true;
        status.last_x = event.clientX;
        status.last_y = event.clientY;
    }
    dom.canvas._mouseup = function () {
        status.draw = false;
    }
    dom.canvas._mouseover = function (event) {
        if(status.draw){
            // console.log(event);
            var x = event.clientX
            var y = event.clientY
            p.beginPath();
            p.moveTo(status.last_x,status.last_y);
            p.lineTo(x,y);
            p.stroke();
            p.closePath();
            status.last_x = event.clientX;
            status.last_y = event.clientY;
        }
    }

    dom.canvas._clear = function(){
        p.clearRect(0,0,dom.canvas.width(),dom.canvas.height());
    }

    dom.close._click = function() {
        dom.group.hide();
        dom.canvas._clear();
    };

    dom.pen._click = function() {
        dom.group.show();
    };
    dom.clear._click = function() {
        dom.canvas._clear();
    };

    dom.group._resize = function(){
        dom.canvas.get(0).width = document.body.scrollWidth;
        dom.canvas.get(0).height = document.body.scrollHeight;
    };

    dom.canvas._init();
    dom.pen.on('click', dom.pen._click);
    dom.close.on('click', dom.close._click);
    dom.clear.on('click', dom.clear._click);
    window.onresize = dom.canvas._init;
    dom.canvas.on('mousemove',dom.canvas._mouseover);
    dom.canvas.on('mousedown',dom.canvas._mousedown);
    dom.canvas.on('mouseup',dom.canvas._mouseup);
})();
