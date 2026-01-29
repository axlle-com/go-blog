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
        $('form[action^="/messages/"]').on('submit', function (evt) {
            evt.preventDefault();
            const form = $(this);
            
            const request = new _glob.request(form);
            request.send(function(response) {
                if (response && response.message) {
                    _glob.noty.success(response.message);
                } else {
                    _glob.noty.success(_glob.t('ui.message.message_sent', 'Message sent successfully'));
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
const _agree = {
    openAgreementModal: function () {
        $('body').on('click', '.js-custom-modal-open', function (e) {
            e.preventDefault();
            const selector = $(this).data('modalName');
            const modal = $(selector);
            if (typeof modal.modal === 'function') {
                modal.modal('show');
                return;
            }
            modal.addClass('modal-agreement');
        });
    },
    closeAgreementModal: function () {
        $('body').on('click', '.js-custom-modal-close', function (e) {
            e.preventDefault();
            const modal = $(this).closest('.js-custom-modal');
            if (typeof modal.modal === 'function') {
                modal.modal('hide');
                return;
            }
            modal.removeClass('modal-agreement');
        });
    },
    run: function () {
        this.openAgreementModal();
        this.closeAgreementModal();
    }
};
const _config = {
    run: function () {
        _message.run();
        _agree.run();
    }
}

$(document).ready(function () {
    _glob.run();
    _config.run();
})
