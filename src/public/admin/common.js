const _auth = {
    login: function () {
        const _this = this;
        $('body').on('click', '.js-submit-button', function () {
            const form = $(this).closest('form');
            const request = new _glob.request(form);
            request.send();
        })
    },
    run: function () {
        $('form').submit(function (evt) {
            evt.preventDefault();
        });
        this.login();
    }
};
const _form = {
    _block: {},
    confirm: function (button, title) {
        const _this = this;
        _this._block.on('click', button, function (e) {
            const saveButton = $(this);
            Swal.fire({
                icon: 'warning',
                title: title,
                text: 'Изменения нельзя будет отменить',
                showDenyButton: true,
                confirmButtonText: 'Сохранить',
                denyButtonText: 'Отменить',
            }).then((result) => {
                if (result.isConfirmed) {
                    _this.send(saveButton);
                } else if (result.isDenied) {
                    Swal.fire('Изменения не сохранены', '', 'info');
                }
            });
        });
    },
    send: function (saveButton) {
        const _this = this;
        const form = saveButton.closest('#global-form');
        const method = form.attr('method');
        if (form) {
            const request = new _glob.request(form).setPreloader('.js-product');
            if (method) {
                request.setMethod(method);
            }
            request.send((response) => {
                let html = $(response.data.view);
                _this._block.html(html);
                _glob.images = {};
                _config.run();
                Swal.fire('Сохранено', '', 'success');
            });
        }
    },
    run: function (selector) {
        this._block = $(selector);
        if (this._block.length) {
            this.confirm('.js-save-button', 'Вы уверены что хотите сохранить все изменения?');
            this.confirm('.js-delete-button', 'Вы уверены что хотите удалить?');
        }
    }
};
const _filter = {
    _block: {},
    send: function () {
        const _this = this;
        _this._block.on('click', '.js-filter-button', function (e) {
            const form = $('#index-form-filter');
            const method = form.attr('method');
            if (form) {
                const request = new _glob.request(form).setPreloader('.js-product');
                if (method) {
                    request.setMethod(method);
                }
                request.send((response) => {
                    let html = $(response.data.view);
                    _this._block.html(html);
                    _glob.images = {};
                    _config.run();
                });
            }
        });

    },
    run: function (selector) {
        this._block = $(selector);
        if (this._block.length) {
            this.send(selector);
        }
    }
};
const _image = {
    confirm: (obj, image) => {
        Swal.fire({
            icon: 'warning',
            title: 'Вы уверены что хотите удалить изображение?',
            text: 'Изменения нельзя будет отменить',
            showDenyButton: true,
            confirmButtonText: 'Удалить',
            denyButtonText: 'Отменить',
        }).then((result) => {
            if (result.isConfirmed) {
                const request = new _glob.request(obj).setMethod('delete').setPreloader('.js-product');
                request.send((response) => {
                    if (response.message) {
                        _glob.noty.success(response.message);
                    } else {
                        _glob.noty.success('Изображение удалено');
                    }
                    image.remove();
                });
            } else if (result.isDenied) {
                Swal.fire('Изображение не удалено', '', 'info');
            }
        })
    },
    add: function () {
        const _this = this;
        $('body').on('change', '.js-image-upload', function () {
            let input = $(this);
            let action = input.attr('data-action');
            let div = input.closest('fieldset').find('.js-image-block');
            let file = input[0].files[0];
            if (!file) {
                return
            }
            let formData = new FormData();
            formData.append('file', file);
            const request = new _glob.request(formData)
                .setPreloader('.js-product')
                .setAction(action);
            request.sendForm((response) => {
                if (response.data.image) {
                    _this.draw(div, response.data.image)
                }
                if (response.message) {
                    _glob.noty.success(response.message);
                } else {
                    _glob.noty.success('Изображение загружено');
                }
            });
        });
    },
    draw: function (div, image) {
        const _this = this;
        if (div.length && image) {
            const imageBlock = `
                    <div class="image-box" style="background-image: url(${image}); background-size: cover;background-position: center;"></div>
                    <div class="overlay-content text-center justify-content-end">
                        <div class="btn-group mb-1" role="group">
                            <a data-fancybox="gallery" href="${image}">
                                <button type="button" class="btn btn-link btn-icon text-danger">
                                    <i class="material-icons">zoom_in</i>
                                </button>
                            </a>
                            <button 
                            type="button" 
                            class="btn btn-link btn-icon text-danger" 
                            data-action="/admin/file/image${image}"
                            data-js-image-delete>
                                <i class="material-icons">delete</i>
                            </button>
                        </div>
                        <input type="hidden" name="image" value="${image}">
                    </div>
                `;
            $(div).html(imageBlock);
            $('.js-image-block-empty').remove();
            _config.fancybox();
        }
    },
    delete: function () {
        const _this = this;
        $('body').on('click', '[data-js-image-delete]', function (evt) {
            let image = $(this).closest('.js-image-block').find('.image-box');
            let input = $(this).closest('.js-image-block').find('input[name="image"]');
            if (!image.length || !input.length) {
                return;
            }
            const action = $(this).attr('data-js-image-href');
            const request = new _glob.request({action: action})
                .setMethod('delete')
                .setPreloader('.js-product');
            request.sendForm((response) => {
                if (response.message) {
                    _glob.noty.success(response.message);
                } else {
                    _glob.noty.success('Изображение удалено');
                }
            });
            input.val('');
            image.remove();
        });
    },
    addArray: function () {
        const _this = this;
        $('body').on('change', '.js-gallery-input', function (evt) {
            let input = $(this);
            let action = input.attr('data-action');
            let files = input[0].files;
            if (!files) {
                return
            }
            let formData = new FormData();
            for (let i = 0; i < files.length; i++) {
                formData.append('files', files[i]);
            }
            new _glob.request(formData).setAction(action).sendForm((response) => {
                if (response.data.images) {
                    let idGallery = input.attr('data-gallery-number');
                    if (!idGallery) {
                        idGallery = _glob.uuid();
                        input.attr('data-gallery-number', idGallery);
                    }
                    _this.drawArray(response.data.images, idGallery)
                }
                if (response.message) {
                    _glob.noty.success(response.message);
                } else {
                    _glob.noty.success('Изображение загружено');
                }
            });
        });
    },
    drawArray: function (images, idGallery) {
        if (Object.keys(images).length <= 0) {
            return
        }
        const _this = this;
        const selector = `[data-gallery-number="${idGallery}"]`;
        const id = $(selector).attr('data-gallery-id');
        const count = $(selector).closest('.js-galleries-general-block').find('.js-gallery-item').length;
        const block = $(selector).closest('.js-galleries-general-block').find('.js-gallery-block-saved');
        for (let i = 0; i < images.length; i++) {
            let number = count + i;
            let url = images[i];
            let image = `<div class="md-block-5 js-gallery-item sort-handle">
                                <div class="img rounded">
                                    <div class="image-box" style="background-image: url(${url}); background-size: cover;background-position: center;"></div>
                                    <div class="overlay-content text-center justify-content-end">
                                        <div class="btn-group mb-1" role="group">
                                            <a data-fancybox="gallery" href="${url}">
                                                <button type="button" class="btn btn-link btn-icon text-danger">
                                                    <i class="material-icons">zoom_in</i>
                                                </button>
                                            </a>
                                            <button type="button" class="btn btn-link btn-icon text-danger" data-js-image-array-id="${idGallery}.${i}">
                                                <i class="material-icons">delete</i>
                                            </button>
                                        </div>
                                    </div>
                                </div>
                                <div>
                                    <input type="hidden" name="galleries[${idGallery}][images][${number}][id]" value="">
                                    <input type="hidden" name="galleries[${idGallery}][images][${number}][gallery_id]" value="${_glob.isEmpty(id) ? '' : id}">
                                    <input type="hidden" name="galleries[${idGallery}][images][${number}][file]" value="${url}">
                                    <div class="form-group small">
                                        <input class="form-control form-shadow" placeholder="Заголовок" name="galleries[${idGallery}][images][${number}][title]" value="">
                                        <div class="invalid-feedback"></div>
                                    </div>
                                    <div class="form-group small">
                                        <input class="form-control form-shadow" placeholder="Описание" name="galleries[${idGallery}][images][${number}][description]" value="">
                                        <div class="invalid-feedback"></div>
                                    </div>
                                    <div class="form-group small">
                                        <input class="form-control form-shadow" type="number" placeholder="Сортировка" name="galleries[${idGallery}][images][${number}][sort]" value="">
                                        <div class="invalid-feedback"></div>
                                    </div>
                                </div>
                            </div>`;
            block.append(image);
        }
        _config.fancybox();
        _config.sort();
    },
    deleteArray: function () {
        const _this = this;
        $('body').on('click', '[data-js-image-array-id]', function (evt) {
            let image = $(this).closest('.js-gallery-item');
            if (!image.length) {
                image = $(this).closest('fieldset').find('.image-box');
                if (!image.length) {
                    return;
                }
            }
            const action = $(this).attr('data-js-image-href');
            let id = $(this).attr('data-js-image-array-id');
            let idGall;
            if (id) {
                let arr = id.split('.');
                idGall = arr[0];
                id = arr[1];
            }
            if (idGall && id) {
                delete _glob.images[idGall]['images'][id];
                image.remove();
            } else {
                _this.confirm({action}, image)
            }
        });
    },
    gallerySort: function () {

    },
    imageBlockEmpty: function () {
        return `
            <div class="form-group js-image-block-empty">
                <label class="control-label button-100" for="js-image-upload">
                    <a type="button" class="btn btn-primary button-image">
                        Загрузить фото
                    </a>
                </label>
                <input
                        type="file"
                        data-action="/admin/file/image"
                        id="js-image-upload"
                        class="custom-input-file js-image-upload"
                        name="file"
                        accept="image/*">
                <div class="invalid-feedback"></div>
            </div>
        `;
    },
    run: function () {
        this.add();
        this.delete();
        this.deleteArray();
        this.addArray();
        this.gallerySort();
    }
};
const _post = {
    login: function () {
        const _this = this;
        $('body').on('click', '.js-submit-button', function () {
            const form = $(this).closest('form');
            const request = new _glob.request(form);
            request.send();
        })
    },
    run: function () {

    }
};
const _infoBlock = {
    add: function () {
        const _this = this;
        $('body').on('click', '.js-info-block-add', function () {
            let div = $(this).closest('.js-info-block-select-form');
            let select = div.find('select').val();
            const request = new _glob.request({action: '/admin/ajax/info-block/get-for-resource/' + select});
            request.setMethod('GET').send((response) => {
                if (response.status) {
                    let html = $(response.data.view);
                    $('.js-info-block-general-block').append(html);
                    _config.select2();
                }
            });
        });
    },
    delete: function () {
        const _this = this;
        $('body').on('click', '.js-info-block-item-delete', function (evt) {
            const id = $(this).attr('data-id');
            const action = $(this).attr('data-href');
            const block = $(this).closest('.js-info-block-item');
            if (id && action) {
                const request = new _glob.request({action});
                request.setMethod('DELETE').send((response) => {
                    if (response.status) {
                        block.remove();
                    }
                });
            } else {
                block.remove();
            }
        });
    },
    run: function () {
        this.add();
        this.delete();
    }
};
const _config = {
    sort: function () {
        let block = document.querySelectorAll('.sortable');
        if (block.length) {
            block.forEach(function (el) {
                const swap = el.classList.contains('swap')
                Sortable.create(el, {
                    swap: swap,
                    animation: 150,
                    handle: '.sort-handle',
                    filter: '.remove-handle',
                    onFilter: function (evt) {
                        evt.item.parentNode.removeChild(evt.item)
                    },
                    onSort: function (evt) {
                        let blocks0 = $(evt.item).closest('.swap').find('[name$="[sort]"]');
                        let blocks1 = $(evt.item).closest('.swap').find('[name$="[property_value_sort]"]');
                        if (blocks0.length) {
                            $.each(blocks0, function (i, value) {
                                $(this).val(i + 1);
                            });
                        }
                        if (blocks1.length) {
                            $.each(blocks1, function (i, value) {
                                $(this).val(i + 1);
                            });
                        }
                    },
                })
            })
        }
    },
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
        if ($('.a-block .sortable').length) {
            this.sort();
        }
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
    _form.run('.a-block-inner');
    _filter.run('.a-block-inner');
    _image.run();
    _infoBlock.run();
    _config.run();
    _auth.run();
    _post.run();
})
