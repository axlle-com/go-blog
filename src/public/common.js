const _config = {
    fancybox: function () {
        try {
            Fancybox.bind('[data-fancybox]', {});
        } catch (e) {
            _glob.console.error(e.message);
        }
    },
    dateRangePicker: function () {
        try {
            flatpickr('.date-range-picker', {
                mode: 'range',
                'locale': 'ru',
                dateFormat: 'd.m.Y',
            });
        } catch (e) {
            _glob.console.error(e.message);
        }
        try {
            flatpickr('.datepicker-wrap', {
                allowInput: true,
                clickOpens: false,
                wrap: true,
                'locale': 'ru',
                dateFormat: 'd.m.Y',
            })
        } catch (e) {
            _glob.console.error(e.message);
        }
    },
    summernote500: function () {
        const summernote500 = $('.summernote-500');
        if (summernote500.length) {
            summernote500.summernote({
                height: 500
            });
        }
    },
    summernote: function () {
        const summernote = $('.summernote');
        if (summernote.length) {
            summernote.summernote({
                height: 150
            });
        }
    },
    flatpickr: function () {
        const selector = '.datetimepicker-inline';
        if ($(selector).length) {
            flatpickr(selector, {
                enableTime: true,
                inline: true
            });
        }
    },
    select2: function () {
        for (const el of document.querySelectorAll('.select2')) {
            let config = {
                width: '100%',
                minimumResultsForSearch: 'Infinity', // hide search
            }
            // live search
            if (el.dataset.select2Search) {
                if (el.dataset.select2Search === 'true') {
                    delete config.minimumResultsForSearch
                }
            }
            // custom content
            if (el.dataset.select2Content) {
                if (el.dataset.select2Content === 'true') {
                    config.templateResult = state => state.id ? $(state.element.dataset.content) : state.text
                    config.templateSelection = state => state.id ? $(state.element.dataset.content) : state.text
                }
            }
            // run
            $(el).select2(config).on('select2:unselecting', function () {
                $(this).data('unselecting', true)
            }).on('select2:opening', function (e) {
                if ($(this).data('unselecting')) {
                    $(this).removeData('unselecting')
                    e.preventDefault()
                }
            })
        }
    },
    run: function () {
        this.select2();
        this.fancybox();
        this.dateRangePicker();
        this.summernote500();
        this.summernote();
        this.flatpickr();
    }
}

$(document).ready(function () {
    _glob.run();
})
