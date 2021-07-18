(function ($) {
    'use strict';
    $(function () {
        var addItem = function(item) {
            var li_class = "";
            var checkbox_checked = "";
            var id_value = ' id="' + item.id + '"';
            if (item.completed) {
                li_class = ' class="completed"';
                checkbox_checked = ' checked=""';
            }
            todoListItem.append("<li " + li_class + id_value + "><div class='form-check'><label class='form-check-label'><input class='checkbox' type='checkbox' " + checkbox_checked + "/>" + item.name + "<i class='input-helper'></i></label></div><i class='remove mdi mdi-close-circle-outline'></i></li>");
        }
        var isBlank = function(string) {
            return string == null || string.trim() === "";
        };
        var session_id = "";
        while (isBlank(session_id)) {
            session_id = prompt("What's your name?");
        }
        
        var todoListItem = $('.todo-list');
        var todoListInput = $('.todo-list-input');
        $('.todo-list-add-btn').on("click", function (event) {
            event.preventDefault();

            var item = $(this).prevAll('.todo-list-input').val();

            if (item) {
                $.ajax({
                    url: "/todos",
                    type: "POST",
                    contentType: "application/json",
                    data: JSON.stringify({
                        name: item,
                        session_id: session_id
                    }),
                    success: addItem
                })
                todoListInput.val("");
            }

        });
        $.ajax({
            url: "/todos",
            type: "MYGET",
            contentType: "application/json",
            data: JSON.stringify({
                session_id: session_id
            }),
            success: function (items) {
                items.sort(function (x, y) {
                    if (x.created_at < y.created_at) {
                        return -1
                    } else if (x.created_at == y.created_at) {
                        return 0
                    } else {
                        return 1
                    }
                })
                items.forEach(e => {
                    addItem(e)
                });
            }
        })

        todoListItem.on('change', '.checkbox', function () {
            var todo_id = parseInt($(this).closest("li").attr('id'));
            var isOn = $(this).attr('checked');
            var $self = $(this)
            $.ajax({
                url: "/todos",
                type: "PUT",
                data: JSON.stringify({
                    id:todo_id
                }),
                dataType:"text",
                contentType: "application/json",
                success: function(e) {
                    console.log(e);
                    if (!e.hasOwnProperty('success')) {
                        if (e.completed) {
                            $self.attr('checked', 'checked');
                        } else {
                            $self.removeAttr('checked');
                        }
                    }
                }
            })

            $(this).closest("li").toggleClass('completed');

        });

        todoListItem.on('click', '.remove', function () {
            var id = $(this).closest("li").attr('id');
            var $self = $(this);
            $.ajax({
                url: "todos/" + id,
                type: "DELETE",
                success: function(e) {
                    if (e.success) {
                        $self.parent().remove();
                    }
                }
            })
        });

    });
})(jQuery);