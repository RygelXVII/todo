var incompleteTasksHolder = $();
var completedTasksHolder = $();

// ---  Function for Deleting an existing task ---
var deleteTask = function (e) {
    var th = $(this);
    $.post('/delete', {id: e.data.id}, function( data ) {
        th.parent().remove();
    });
};

// ---  Function for Adding a new task ---
var addTask = function() {
    var taskInput = $('#new-task');
    
    if ( taskInput.val() == '' ) {
        alert('You Can\'t add a Empty Task');
    } else {
        let valInput = taskInput.val();
        $.post('/create', {task: valInput}, function( data ) {
            let listItem = $('<li><input type="checkbox"><label>' + valInput + '</label><input type="text"><button class="delete">Delete</button></li>');
            incompleteTasksHolder.append(listItem);
            taskInput.val('');
            //Bind Event Handler to the task that was Added
            let newTaskAdded = incompleteTasksHolder.find('li').last();
            newTaskAdded.find('.delete').on('click', {id: data}, deleteTask);
            newTaskAdded.find('input[type="checkbox"]').on('click', markTask);
        }, 'text')
    }
};

// ---  Function for Marking a task Complete or Incomplete ---
var markTask = function() {
    //This function will check if the task is in edit mode, if true will change it
    function checkEditState(task) {
        if ( task.hasClass('editMode') ) {
            task.find('label').text( task.find('input[type="text"]').val() )
            task.find('input[type="text"]').removeAttr('value');
            task.removeAttr('class');
        }
    }

    if ( $(this).closest('ul').attr('id') == 'incomplete-tasks' ) {
        var completedTask = $(this).parent().remove();
        checkEditState(completedTask);
        completedTask.find('input[type="checkbox"]').prop( "checked", true );
        bindTaskEvents( completedTask );
        completedTasksHolder.append(completedTask);
    } else {
        var taskToComplete = $(this).parent().remove();
        checkEditState(taskToComplete);
        taskToComplete.find('input[type="checkbox"]').removeAttr('checked');
        bindTaskEvents( taskToComplete );
        incompleteTasksHolder.append(taskToComplete);
    }
};

// --- Function for binding Event Handler ---
function bindTaskEvents (item) {
    item.find('.delete').on('click', deleteTask);
    item.find('input[type="checkbox"]').on('click', markTask);
}

$(document).ready(function(){
    //--- The TodoList in JQuery ---
    incompleteTasksHolder = $('#incomplete-tasks');
    completedTasksHolder = $('#completed-tasks');
    // Binding Event handlers to Add, checkboxes, edit and delete
    let buttonAddTask = $('#new-task').next();
    buttonAddTask.on('click', addTask);
    bindTaskEvents( $('li') );
})