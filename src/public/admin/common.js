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
    confirm: function (button, title, confirmButtonText = 'Save') {
        const _this = this;
        _this._block.on('click', button, function (e) {
            const saveButton = $(this);
            Swal.fire({
                icon: 'warning',
                title: title,
                text: 'You will not be able to undo this action',
                showDenyButton: true,
                confirmButtonText: confirmButtonText,
                denyButtonText: 'Cancel',
            }).then((result) => {
                if (result.isConfirmed) {
                    _this.send(saveButton);
                } else if (result.isDenied) {
                    Swal.fire('Changes not saved', '', 'info');
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
                Swal.fire('Saved', '', 'success');
            });
        }
    },
    run: function (selector) {
        this._block = $(selector);
        if (this._block.length) {
            this.confirm('.js-save-button', 'Are you sure you want to save all changes?');
            this.confirm('.js-delete-button', 'Are you sure you want to delete?', 'Delete');
        }
    }
};
const _filterApi = {
    create: function () {
        return {
            // target: '#id' | '.class' | DOM | jQuery
            send: function (target) {
                const $list = (target && target.jquery) ? target : $(target);
                if (!$list.length) {
                    console.warn(`_filterApi: no elements found for selector "${target}"`);
                    return;
                }

                $list.each(function () {
                    const $el = $(this);

                    // custom URL for each select
                    const getUrl = () => $el.data('action') || $el.attr('data-action');
                    const url = getUrl();
                    if (!url) {
                        console.warn('select2: data-action not set for', $el[0]);
                        return;
                    }

                    // if already initialized — recreate
                    if ($el.hasClass('select2-hidden-accessible')) {
                        $el.select2('destroy');
                    }

                    $el.select2({
                        placeholder: $el.data('placeholder') || 'Select...',
                        allowClear: true,
                        minimumInputLength: 0,
                        width: '100%',
                        ajax: {
                            url: getUrl,
                            dataType: 'json',
                            delay: 250,
                            data: function (params) {
                                const q = params.term || '';
                                return q ? {query: q, page: params.page || 1} : {page: params.page || 1};
                            },
                            processResults: function (resp, params) {
                                params.page = params.page || 1;

                                const data = resp.data || resp || {};
                                const items = Array.isArray(data.items) ? data.items : [];
                                const p = data.paginator || {};
                                const pageSize = Number(p.pageSize || 20);
                                const page = Number(p.page || params.page || 1);
                                const total = Number(p.total || 0);

                                const results = items.map(it => ({
                                    id: it.id ?? it.ID ?? it.value,
                                    text: it.text ?? it.title ?? it.Name
                                })).filter(o => o.id != null && o.text != null);

                                const more = page * pageSize < total;
                                return {results, pagination: {more}};
                            },
                            cache: true
                        }
                    });
                });
            }
        };
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
            title: 'Are you sure you want to delete the image?',
            text: 'You will not be able to undo this action',
            showDenyButton: true,
            confirmButtonText: 'Delete',
            denyButtonText: 'Cancel',
        }).then((result) => {
            if (result.isConfirmed) {
                const request = new _glob.request(obj).setMethod('delete').setPreloader('.js-product');
                request.send((response) => {
                    if (response.message) {
                        _glob.noty.success(response.message);
                    } else {
                        _glob.noty.success('Image deleted');
                    }
                    image.remove();
                });
            } else if (result.isDenied) {
                Swal.fire('Image not deleted', '', 'info');
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
            const resource = $('body').find('[name="resource"]');
            let formData = new FormData();
            formData.append('file', file);

            if (resource[0] && $(resource[0]).val()) {
                formData.append('resource', $(resource[0]).val());
            }

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
                    _glob.noty.success('Image uploaded');
                }
            });
        });
    },
    draw: function (div, image) {
        const _this = this;
        if (div.length && image) {
            $(div).html(_this.imageBlock(image));
            _config.fancybox();
        }
    },
    delete: function () {
        const _this = this;
        $('body').on('click', '[data-js-image-delete]', function (evt) {
            let block = $(this).closest('.js-image-block');
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
                    _glob.noty.success('Image deleted');
                }
            });
            input.val('');
            block.html('');
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
            formData.append('resource', 'galleries');
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
                    _glob.noty.success('Image uploaded');
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
                                            <button type="button" 
                                                class="btn btn-link btn-icon text-danger" 
                                                data-js-gallery-image-href="/admin/file/image${url}"
                                                data-js-gallery-image-delete>
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
                                        <input class="form-control form-shadow" placeholder="Title" name="galleries[${idGallery}][images][${number}][title]" value="">
                                        <div class="invalid-feedback"></div>
                                    </div>
                                    <div class="form-group small">
                                        <input class="form-control form-shadow" placeholder="Description" name="galleries[${idGallery}][images][${number}][description]" value="">
                                        <div class="invalid-feedback"></div>
                                    </div>
                                    <div class="form-group small">
                                        <input class="form-control form-shadow" type="number" placeholder="Sort" name="galleries[${idGallery}][images][${number}][sort]" value="">
                                        <div class="invalid-feedback"></div>
                                    </div>
                                </div>
                            </div>`;
            block.append(image);
        }
        _config.fancybox();
        _config.sort();
    },
    deleteInGallery: function () {
        const _this = this;
        $('body').on('click', '[data-js-gallery-image-delete]', function (evt) {
            let image = $(this).closest('.js-gallery-item');
            const action = $(this).attr('data-js-gallery-image-href');
            const request = new _glob.request({action: action})
                .setMethod('delete')
                .setPreloader('.js-product');
            request.sendForm((response) => {
                if (response.message) {
                    _glob.noty.success(response.message);
                } else {
                    _glob.noty.success('Image deleted');
                }
            });
            image.remove();
        });
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
    imageBlock: function (image) {
        return `
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
                        data-js-image-href="/admin/file/image${image}"
                        data-js-image-delete>
                        <i class="material-icons">delete</i>
                    </button>
                </div>
                <input type="hidden" name="image" value="${image}">
            </div>
        `;
    },
    imageBlockEmpty: function () {
        return `
            <div class="form-group js-image-block-empty">
                <label class="control-label button-100" for="js-image-upload">
                    <a type="button" class="btn btn-primary button-image">
                        Upload photo
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
        this.deleteInGallery();
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
    _block: {},
    add: function () {
        const _this = this;
        $('body').on('click', '.js-info-blocks-add', function (evt) {
            _this._block = $('.js-info-block-saved');
            const select = $(this).closest('.js-info-blocks-general-block').find('.js-info-blocks-select');
            const action = select.find('option:selected').data('action');
            if (!action) {
                _glob.console.error('Empty identifier')
                return
            }
            const request = new _glob.request({action});
            request.setMethod('GET').send((response) => {
                select.val(null).trigger('change');
                let html = $(response.data.view);
                const count = $('input[name^="info_blocks["][name$="[id]"]').length;
                html.find('[name^="info_blocks["]').each(function () {
                    const $el = $(this);
                    const name = $el.attr('name');
                    const newName = name.replace(/^info_blocks\[\d+\]/, 'info_blocks[' + count + ']');
                    $el.attr('name', newName);
                });
                _this._block.append(html);
                _config.run();
            });
        });
    },
    delete: function () {
        const _this = this;
        $('body').on('click', '.js-info-block-item-delete', function (evt) {
            const id = $(this).attr('data-id');
            const action = $(this).attr('data-action');
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
    detach: function () {
        const _this = this;
        $('body').on('click', '.js-info-blocks-detach', function (evt) {
            const action = $(this).attr('data-action');
            const block = $(this).closest('.js-info-blocks-item');
            if (action) {
                const request = new _glob.request({action});
                request.setMethod('DELETE').send((response) => {
                    block.remove();
                });
            } else {
                block.remove();
            }
        });
    },
    run: function () {
        this._block = $('.js-info-block-saved');
        this.add();
        this.delete();
        this.detach();
    }
};
const _menu = {
    _block: {},
    search: function () {
        const _this = this;
        const filterApi = _filterApi.create();
        filterApi.send('.select2-search')
    },
    setUrl: function () {
        const _this = this;
        const sel = 'select[name^="menu_items"][name$="[publisher_uuid]"]';

        // remove old handlers and attach new ones
        $(document)
            .off('select2:select.menuurl select2:clear.menuurl', sel)
            .on('select2:select.menuurl', sel, function (e) {
                const select = e.target;
                const url = select.options[select.selectedIndex]?.dataset.url || '';

                const fieldset = select.closest('.form-block.js-menu-items-publisher-url');
                if (!fieldset) return;

                const linkInput = fieldset.querySelector('input[name$="[url]"]');
                if (!linkInput) return;

                if (!linkInput.dataset.oldUrl) {
                    linkInput.dataset.oldUrl = linkInput.value; // remember original
                }
                linkInput.value = url;
            })
            .on('select2:clear.menuurl', sel, function (e) {
                const select = e.target;
                const fieldset = select.closest('.form-block.js-menu-items-publisher-url');
                if (!fieldset) return;

                const linkInput = fieldset.querySelector('input[name$="[url]"]');
                if (!linkInput) return;

                linkInput.value = linkInput.dataset.oldUrl || '';
                delete linkInput.dataset.oldUrl;
            });
    },
    run: function (selector) {
        this._block = $(selector);
        if (this._block.length) {
            this.setUrl();
            this.search();
        }
    }
};
const _template = {
    _block: {},
    add: function () {
        const _this = this;
        $('body').on('change', '.js-template-select', function (evt) {
            const action = $(this).val();
            const block = $(this).closest('.a-block-inner').find('#HTML');
            if (action) {
                const request = new _glob.request({action});
                request.setMethod('GET').send((response) => {
                    block.html('')
                    MyCodeMirror6.createEditor(document.getElementById('HTML'), response.data.view);
                });
            }
        });
    },
    run: function () {
        const htmlInit = $('#HTML_container').html();
        $('#HTML_container').html('');
        const jsInit = $('#JS_container').html();
        $('#JS_container').html('');
        const cssInit = $('#CSS_container').html();
        $('#CSS_container').html('');

        MyCodeMirror6.createEditor(document.getElementById('HTML_container'), htmlInit);
        MyCodeMirror6.createEditor(document.getElementById('JS_container'), jsInit);
        MyCodeMirror6.createEditor(document.getElementById('CSS_container'), cssInit);

        // before submit copy values to textarea
        $('body').on('click', '.js-save-button', () => {
            $('#HTML').val(document.getElementById('HTML_container').getValue());
            $('#JS').val(document.getElementById('JS_container').getValue());
            $('#CSS').val(document.getElementById('CSS_container').getValue());
            $('#global-form').trigger('submit');
        });

        this.add();
    }
};
const _message = {
    _selector: '.js-message-content',
    _selectorList: '.js-message-list',
    _unviewed: '.js-message-unviewed',
    read: function () {
        const _this = this;
        $('body').on('click', '.mail-item', function (e) {
            const $this = $(this);
            const action = $this.data('jsMessageAction');
            if (action) {
                $(_this._selector).html('')
                const request = new _glob.request({action});
                request.setMethod('GET').send((response) => {
                    $(_this._selector).html(response.data.view)
                    $(_this._selectorList).html(response.data.list)
                    $(_this._unviewed).html(response.data.unviewed)
                });
            }
        });
    },
    delete: function () {
        const _this = this;
        $('body').on('click', '.js-message-content-delete', function (e) {
            const $this = $(this);
            const action = $this.data('jsMessageAction');
            if (action) {
                $(_this._selector).html('')
                const request = new _glob.request({action});
                request.setMethod('DELETE').send((response) => {
                    $(_this._selectorList).html(response.data.view)
                });
            }
        });
    },
    run: function () {
        this.read();
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
            if (el.classList.contains('select2-search')) {
                continue;
            }

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
        _menu.run('.a-block-inner.menu');
    }
}

$(document).ready(function () {
    _glob.run();
    _config.run();
    _form.run('.a-block-inner');
    _filter.run('.a-block-inner');
    _image.run();
    _infoBlock.run();
    _auth.run();
    _post.run();
    _template.run();
    _message.run();
    _menu.run('.a-block-inner.menu');
})
