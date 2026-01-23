const _message = {
    login: function () {
        const _this = this;
        $('body').on('click', '.js-submit-button', function () {
            const form = $(this).closest('form');
            const request = new _glob.request(form);
            request.send();
        })
    },
    contact: function () {
        const _this = this;
        $('form[action="/messages"]').on('submit', function (evt) {
            evt.preventDefault();
            const form = $(this);
            
            const request = new _glob.request(form);
            request.send(function(response) {
                if (response && response.message) {
                    _glob.noty.success(response.message);
                } else {
                    _glob.noty.success(_glob.t('ui.success.message_sent', 'Message sent successfully'));
                }
            });
        });
    },
    run: function () {
        $('form').submit(function (evt) {
            evt.preventDefault();
        });
        this.login();
        this.contact();
    }
};
const _config = {
    run: function () {
        _message.run();
    }
}

$(document).ready(function () {
    _glob.run();
    _config.run();
})
