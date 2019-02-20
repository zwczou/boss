$(function() {
  var validator = $('form').validate({
    ignore: 'input[type=hidden], .select2-search__field', // ignore hidden fields
    errorClass: 'validation-invalid-label',
    successClass: 'validation-valid-label',
    validClass: 'validation-valid-label',
    highlight: function(element, errorClass) {
      $(element).removeClass(errorClass);
    },
    unhighlight: function(element, errorClass) {
      $(element).removeClass(errorClass);
    },
    success: function(label) {
      label.addClass('validation-valid-label').text('验证通过'); // remove to hide Success message
    },

    // Different components require proper error label placement
    errorPlacement: function(error, element) {

      // Unstyled checkboxes, radios
      if (element.parents().hasClass('form-check')) {
        error.appendTo( element.parents('.form-check').parent() );
      }

      // Input with icons and Select2
      else if (element.parents().hasClass('form-group-feedback') || element.hasClass('select2-hidden-accessible')) {
        error.appendTo( element.parent() );
      }

      // Input group, styled file input
      else if (element.parent().is('.uniform-uploader, .uniform-select') || element.parents().hasClass('input-group')) {
        error.appendTo( element.parent().parent() );
      }

      // Other elements
      else {
        error.insertAfter(element);
      }
    },
    rules: {
      username: {
        required: true,
        minlength: 5
      },
      password: {
        required: true,
        minlength: 6
      },
    },
    messages: {
    }
  });
});
