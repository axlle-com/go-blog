const _cl_ = (p) => {
    console.log(p);
}
const _glob = {
    ERROR_MESSAGE: 'Произошла ошибка, попробуйте позднее!',
    ERROR_FIELD: 'Поле обязательное для заполнения',
    spareParts: [],
    images: {},
    path: null,
    pathArray: null,
    pathSearchParams: null,
    pathHash: null,
    console: {
        error: function (message = null) {
            if (message) {
                console.log(`%c ${message} `, `background: #d43f3a; color: #eee`);
            } else {
                console.log(`%c ${_glob.ERROR_MESSAGE} `, `background: #d43f3a; color: #eee`);
            }
        },
        info: function (message) {
            console.log(`%c ${message} `, `background: #4cae4c; color: #eee`);
        },
    },
    noty: {
        config: function (type, message) {
            if (typeof Noty !== 'undefined') {
                const text = '<h5>Внимание</h5>' + message;
                const _config = {type, text, timeout: 4000, theme: 'relax'};
                new Noty(_config).show();
            } else {
                _glob.console.error(message);
                alert(message);
            }
        },
        error: function (message = 'Произошла ошибка!') {
            this.config('error', message);
        },
        success: function (message = 'Все прошло успешно!') {
            this.config('success', message);
        },
        info: function (message = 'Обратите внимание!') {
            this.config('info', message);
        }
    },
    preloader: {
        block: `<div class="preloader" style="display: none;"><div class="lds-spinner"><div></div><div></div><div></div><div></div><div></div><div></div><div></div><div></div><div></div><div></div><div></div><div></div></div></div>`,
        style: `<style id="preloader-style">.relative{position:relative}.preloader{position:absolute;top:0;left:0;bottom:0;right:0;background:rgba(255,255,255,.30);z-index:10100}.lds-spinner{color:#006400;display:inline-block;width:64px;height:64px;position:absolute;top:10%;left:50%;margin-right:-50%;transform:translate(-50%,-50%)}.lds-spinner div{transform-origin:32px 32px;animation:lds-spinner 1.2s linear infinite}.lds-spinner div:after{content:" ";display:block;position:absolute;top:3px;left:29px;width:5px;height:14px;border-radius:20%;background:#3d8bfd}.lds-spinner div:nth-child(1){transform:rotate(0);animation-delay:-1.1s}.lds-spinner div:nth-child(2){transform:rotate(30deg);animation-delay:-1s}.lds-spinner div:nth-child(3){transform:rotate(60deg);animation-delay:-.9s}.lds-spinner div:nth-child(4){transform:rotate(90deg);animation-delay:-.8s}.lds-spinner div:nth-child(5){transform:rotate(120deg);animation-delay:-.7s}.lds-spinner div:nth-child(6){transform:rotate(150deg);animation-delay:-.6s}.lds-spinner div:nth-child(7){transform:rotate(180deg);animation-delay:-.5s}.lds-spinner div:nth-child(8){transform:rotate(210deg);animation-delay:-.4s}.lds-spinner div:nth-child(9){transform:rotate(240deg);animation-delay:-.3s}.lds-spinner div:nth-child(10){transform:rotate(270deg);animation-delay:-.2s}.lds-spinner div:nth-child(11){transform:rotate(300deg);animation-delay:-.1s}.lds-spinner div:nth-child(12){transform:rotate(330deg);animation-delay:0s}@keyframes lds-spinner{0%{opacity:1}100%{opacity:0}}</style>`,
    },
    request: class {
        validate;
        hasErrors = false;
        hasSend = false;
        payload;
        action;
        form = null;
        response;
        data;
        view;
        preloader;
        method = 'POST';
        isFormReset = true;

        constructor(object = null, validate = true) {
            this.reset().validate = validate;
            if (object) {
                this.setObject(object);
            }
        }

        reset() {
            this.hasErrors = this.hasSend = false;
            this.payload = this.action = this.form = this.response = this.data = this.view = this.preloader = null;
            return this;
        }

        setPreloader(element, top = 10) {
            const _this = this;
            const block = $(element);
            if (block && block.length) {
                const head = $('head');
                const style = $('style#preloader-style');
                if (!style.length) {
                    head.append(_glob.preloader.style);
                }
                _this.preloader = $(_glob.preloader.block);
                block.addClass('relative');
                if (top !== 10) {
                    const ldsSpinner = `<style id="lds-spinner">.lds-spinner{top:${top}%}</style>`;
                    const styleSpinner = $('style#lds-spinner');
                    if (!styleSpinner.length) {
                        head.append(ldsSpinner);
                    } else {
                        styleSpinner.html(`.lds-spinner{top: ${top}%;}`);
                    }
                }
                block.prepend(_this.preloader);
            }
            return _this;
        }

        setObject(object = null) {
            this.data = this.view = this.response = null;
            if (object) {
                if ('action' in object) {
                    this.action = object.action;
                    delete object.action;
                    let data = new FormData();
                    if (Object.keys(object).length) {
                        for (let key in object) {
                            data.append(key, object[key]);
                        }
                    }
                    this.payload = data;
                } else if (object instanceof jQuery) {
                    this.form = object;
                    this.action = this.form.attr('action');
                    this.payload = new FormData(this.form[0]);
                } else if (object instanceof FormData) {
                    this.payload = object;
                } else {
                    _glob.console.error('Не известные данные');
                }
            }
            return this;
        }

        validateForm() {
            if (this.form && this.validate) {
                let err = [];
                $.each(this.form.find('[data-validator-required]'), function (index, value) {
                    err.push(_glob.validation.change($(this)));
                });
                this.hasErrors = err.indexOf(true) !== -1;
            }
        }

        setAction(action) {
            this.action = action;
            return this;
        }

        setMethod(method) {
            this.method = method;
            return this;
        }

        deepSetWithTypes(target, path, value) {
            const isMulti = /\[\]$/.test(path);
            const keys = path.match(/([^\[\]\.]+)/g) || [];
            let curr = target;

            // создаём контейнеры на пути
            for (let i = 0; i < keys.length - 1; i++) {
                const key = keys[i];
                const next = keys[i + 1];
                const needArray = /^\d+$/.test(next);

                const exists = curr[key];
                if (exists == null) {
                    curr[key] = needArray ? [] : {};
                } else {
                    if (needArray && !Array.isArray(exists)) {
                        curr[key] = [];
                    }
                    if (!needArray && (typeof exists !== 'object' || Array.isArray(exists))) {
                        curr[key] = {};
                    }
                }
                curr = curr[key];
            }

            const lastKey = keys[keys.length - 1];
            const typeHint = _glob.typeHints.get(path);

            if (isMulti) {
                if (!Array.isArray(curr[lastKey])) {
                    curr[lastKey] = curr[lastKey] == null ? [] : [curr[lastKey]];
                }
                const incoming = Array.isArray(value) ? value : [value];
                incoming.forEach(item => {
                    const casted = this.castByType(item, typeHint);
                    if (casted !== '' && casted != null) {
                        curr[lastKey].push(casted);
                    }
                });
                return;
            }

            let newValue = this.castByType(value, typeHint);
            if (newValue === '' || newValue == null || (this.isPlainObject(newValue) && Object.keys(newValue).length === 0)) {
                return;
            }

            // если уже было значение и это не массив — превратим в массив и добавим
            if (Object.prototype.hasOwnProperty.call(curr, lastKey) && !Array.isArray(curr[lastKey])) {
                curr[lastKey] = [curr[lastKey]];
            }
            if (Array.isArray(curr[lastKey])) {
                curr[lastKey].push(newValue);
            } else {
                curr[lastKey] = newValue;
            }
        }

        send(callback = null) {
            const _this = this;
            this.validateForm();

            if (this.hasErrors) {
                _glob.noty.error('Заполнены не все обязательные поля');
                return;
            }

            if (this.hasSend) {
                _glob.console.error('Форма еще отправляется');
                return;
            }

            if (this.preloader) {
                this.preloader.show();
            }

            this.hasSend = true;
            let formObject = {};
            const csrf = $('meta[name="csrf-token"]').attr('content');

            // Собираем объект на основе payload
            _this.payload.forEach(function (value, key) {
                _this.deepSetWithTypes(formObject, key, value);
            });

            let ajaxOptions = {
                url: _this.action,
                headers: {'X-CSRF-TOKEN': csrf},
                type: _this.method,
                dataType: 'json',
                beforeSend: function () {
                },
                success: function (response) {
                    _this.setData(response).defaultBehavior();
                    if (callback) {
                        callback(response);
                    }
                },
                error: function (response) {
                    _this.errorResponse(response);
                },
                complete: function () {
                    _this.hasSend = false;
                    if (_this.preloader) {
                        _this.preloader.hide();
                    }
                }
            };

            // Если GET, передаем как обычные query-параметры
            if (_this.method.toUpperCase() === 'GET') {
                ajaxOptions.data = formObject;
            } else {
                // Для POST (или других типов, где нужно тело)
                ajaxOptions.data = JSON.stringify(formObject);
                ajaxOptions.contentType = 'application/json';
            }

            $.ajax(ajaxOptions);
        }

        sendForm(callback = null) {
            const _this = this;
            this.validateForm();
            if (this.hasErrors) {
                _glob.noty.error('Заполнены не все обязательные поля');
                return;
            }
            if (this.hasSend) {
                _glob.console.error('Форма еще отправляется');
                return;
            }
            if (this.preloader) {
                this.preloader.show();
            }
            this.hasSend = true;
            // this.appendImages();
            const csrf = $('meta[name="csrf-token"]').attr('content');
            $.ajax({
                url: _this.action,
                headers: {'X-CSRF-TOKEN': csrf},
                type: _this.method,
                dataType: 'json',
                data: _this.payload,
                processData: false,
                contentType: false,
                beforeSend: function () {
                },
                success: function (response) {
                    _this.setData(response).defaultBehavior();
                    if (!!callback) {
                        callback(response);
                    }
                },
                error: function (response) {
                    _this.errorResponse(response);
                },
                complete: function () {
                    _this.hasSend = false;
                    if (_this.preloader) {
                        _this.preloader.hide();
                    }
                }
            });
        }

        getData() {
            if (!this.data) {
                if (this.response && this.response.status && this.response.data) {
                    this.data = this.response.data;
                } else {
                    this.data = null;
                }
            }
            return this.data;
        }

        setData(response) {
            this.response = response;
            this.data = response.data;
            this.form ? this.form[0].reset() : null;
            return this;
        }

        defaultBehavior() {
            let data, url, redirect, view;
            if ((data = this.getData())) {
                if ((url = data.url)) {
                    this.setLocation(url);
                }
                if ((redirect = data.redirect)) {
                    window.location.href = redirect;
                }
                if ((view = data.view)) {
                    this.view = view;
                }
                try {
                    this.form[0].reset();
                } catch (e) {
                }
            }
        }

        errorResponse(response, form = null) {
            let extraErrors = '';
            let json, message, errors;
            if (response && (json = response.responseJSON)) {
                message = json.message;
                if (message) {
                    message = message.replace(/\|/gi, `<br>`);
                }
                errors = json.errors;
            }

            if (!message && response.responseText) {
                try {
                    message = JSON.parse(response.responseText).message;
                } catch (e) {
                    _glob.console.error(e)
                }
            }

            if (
                response.status === 400
                || response.status === 419
                || response.status === 422
            ) {
                let errs = errors || response.errors;
                if (errs && Object.keys(errs).length) {
                    for (let key in errs) {
                        // подсветка
                        let selector = `[data-validator="${key}"]`;
                        if (form) {
                            $(form).find(selector).addClass('is-invalid');
                        } else {
                            $(selector).addClass('is-invalid');
                        }

                        // собираем текст ошибок
                        let fieldErrors = Array.isArray(errs[key]) ? errs[key] : [errs[key]];
                        extraErrors += fieldErrors.join('<br>') + '<br>';
                    }
                }
                _glob.noty.error(
                    (message ? message + '<br>' : '') + (extraErrors || _glob.ERROR_MESSAGE)
                );
            } else if (response.status === 406) {
                _glob.noty.error(message ? message : _glob.ERROR_MESSAGE);
            } else if (response.status === 500) {
                _glob.noty.error(message ? message : response.statusText);
            } else {
                _glob.noty.error(message ? message : response.statusText);
            }
        }

        setLocation(curLoc) {
            let url = '';
            try {
                url = curLoc === 'index' ? '/' : curLoc;
                history.pushState(null, null, url);
                return;
            } catch (e) {
                _glob.console.error(e.message);
            }
            location.hash = '#' + url;
        }

        castValueAuto(value) {
            // Бинарные/особые типы не трогаем
            if (value instanceof File || value instanceof Blob || value instanceof Date || value instanceof ArrayBuffer) return value;
            if (typeof value !== 'string') return value;

            const s = value.trim();
            if (s === '') return s;

            // boolean
            if (/^(true|false)$/i.test(s)) return s.toLowerCase() === 'true';

            // null/undefined
            if (/^(null)$/i.test(s)) return null;
            if (/^(undefined)$/i.test(s)) return undefined;

            // int
            if (/^[+-]?\d+$/.test(s)) {
                const n = Number(s);
                if (Number.isSafeInteger(n)) return n;
            }
            // float
            if (/^[+-]?\d*\.\d+$/.test(s)) {
                const n = Number(s);
                if (!Number.isNaN(n)) return n;
            }
            return value;
        }

        castByType(value, typeStr) {
            if (!typeStr) {
                return this.castValueAuto(value);
            }

            const isArrayType = /\[\]$/.test(typeStr);
            const base = typeStr.replace(/\[\]$/, '').toLowerCase();

            const castOne = (val) => {
                if (val instanceof File || val instanceof Blob || val instanceof Date || val instanceof ArrayBuffer) return val;
                if (typeof val !== 'string') return val;
                const s = val.trim();
                if (s === '') return s;

                switch (base) {
                    case 'int':
                    case 'integer':
                        if (/^[+-]?\d+$/.test(s)) {
                            const n = Number(s);
                            if (Number.isSafeInteger(n)) {
                                return n;
                            }
                        }
                        return s;
                    case 'float':
                    case 'number':
                        if (/^[+-]?\d+(\.\d+)?$/.test(s)) {
                            const n = Number(s);
                            if (!Number.isNaN(n)) {
                                return n;
                            }
                        }
                        return s;
                    case 'bool':
                    case 'boolean':
                        // поддерживаем "1"/"0" и true/false
                        if (value === '1' || value === 1 || value === true || value === 'true') {
                            return true;
                        }
                        if (value === '0' || value === 0 || value === false || value === 'false') {
                            return false;
                        }
                        return Boolean(value);
                    case 'json':
                        try {
                            return JSON.parse(s);
                        } catch {
                            return s;
                        }
                    case 'string':
                    default:
                        return s;
                }
            };

            if (isArrayType) {
                if (Array.isArray(value)) return value.map(castOne);
                return [castOne(value)];
            }
            return castOne(value);
        }

        isPlainObject(value) {
            if (value === null || typeof value !== 'object') return false;
            const proto = Object.getPrototypeOf(value);
            return proto === Object.prototype || proto === null;
        }
    },
    validation: {
        control: function () {
            const self = this;
            $('body').on('blur', '[data-validator-required]', function () {
                let field = $(this);
                self.change(field);
            })
        },
        change: function (field) {
            let err = false, self = this;
            let help = field.closest('div').find('.invalid-feedback');
            if (field.attr('type') === 'checkbox') {
                if (field.prop('checked')) {
                    field.removeClass('is-invalid');
                    help.text('').hide();
                } else {
                    field.addClass('is-invalid');
                    help.text(_glob.ERROR_FIELD).show();
                    err = true;
                }
            } else {
                if (field.val()) {
                    field.removeClass('is-invalid');
                    help.text('').hide();
                } else {
                    field.addClass('is-invalid');
                    help.text(_glob.ERROR_FIELD).show();
                    err = true;
                }
            }
            return err;
        }
    },
    cookie: class {
        constructor(name, value, options) {
            this.name = name;
            this.value = value;
            this.options = options;
        }

        get() {
            let matches = document.cookie.match(
                new RegExp("(?:^|; )" + this.name.replace(/([\.$?*|{}\(\)\[\]\\\/\+^])/g, '\\$1') + "=([^;]*)"));
            return matches ? decodeURIComponent(matches[1]) : undefined;
        }

        set() {
            this.options = this.options || {};
            let expires = this.options.expires;
            if (typeof expires == "number" && expires) {
                let d = new Date();
                d.setDate(d.getDate() + expires);
                expires = this.options.expires = d;
            }
            if (expires && expires.toUTCString) {
                this.options.expires = expires.toUTCString();
            }
            this.value = encodeURIComponent(this.value);
            let updatedCookie = this.name + "=" + this.value;
            for (let propName in this.options) {
                updatedCookie += "; " + propName;
                let propValue = this.options[propName];
                if (propValue !== true) {
                    updatedCookie += "=" + propValue;
                }
            }
            document.cookie = updatedCookie;
            return this;
        }
    },
    setMaps: function () {
        const cookie = new this.cookie('_maps');
        if (!cookie.get()) {
            cookie.value = true;
            cookie.options = {expires: '', path: '/'};
            cookie.set();
        }
    },
    typeHints: {
        _exact: new Map(),   // полный путь (если когда-нибудь захочешь)
        _norm: new Map(),   // нормализованный путь ( [0] -> [] )
        _leaf: new Map(),   // тип по имени поля (leaf), например parent_id -> int

        _normPath(name) {
            return String(name).replace(/\[\d+\]/g, '[]');
        },
        _lastField(name) {
            const parts = String(name).match(/([^\[\]\.]+)/g) || [];
            return parts[parts.length - 1] || '';
        },

        // задаём типы пачкой (ключи — имена полей или точные пути)
        set(dict) {
            Object.entries(dict || {}).forEach(([key, type]) => {
                const t = String(type).trim();
                if (/[\[\]\.]/.test(key)) {
                    // это путь
                    const p = String(key).trim();
                    this._exact.set(p, t);
                    this._norm.set(this._normPath(p), t);
                } else {
                    // это просто имя поля
                    this._leaf.set(String(key).trim(), t);
                }
            });
            return this;
        },

        // (не обязательно, но пусть будет)
        setPath(path, type) {
            const p = String(path).trim();
            const t = String(type).trim();
            this._exact.set(p, t);
            this._norm.set(this._normPath(p), t);
            return this;
        },

        setLeaf(field, type) {
            this._leaf.set(String(field).trim(), String(type).trim());
            return this;
        },

        // резолв типа
        get(path) {
            const exact = this._exact.get(path);
            if (exact) return exact;

            const norm = this._norm.get(this._normPath(path));
            if (norm) return norm;

            return this._leaf.get(this._lastField(path)) || null;
        },

        run: function () {
            this.set({
                // базовые
                id: 'int',
                uuid: 'string',

                // users
                first_name: 'string',
                last_name: 'string',
                patronymic: 'string',
                phone: 'string',
                email: 'string',
                is_email: 'bool',
                is_phone: 'bool',
                status: 'int',
                avatar: 'string',
                password_hash: 'string',
                remember_token: 'string',
                auth_token: 'string',
                auth_key: 'string',
                password_reset_token: 'string',
                created_at: 'string',
                updated_at: 'string',
                deleted_at: 'string',

                // связи и ACL
                permission_id: 'int',
                user_id: 'int',
                role_id: 'int',
                user_uuid: 'string',
                relation_uuid: 'string',

                // общие для контента/шаблонов
                template_id: 'int',
                post_category_id: 'int',
                meta_title: 'string',
                meta_description: 'string',
                alias: 'string',
                url: 'string',
                is_main: 'bool',
                is_published: 'bool',
                is_favourites: 'bool',
                has_comments: 'bool',
                show_image_post: 'bool',
                show_image_category: 'bool',
                in_sitemap: 'bool',
                media: 'string',
                title: 'string',
                title_short: 'string',
                description_preview: 'string',
                description: 'string',
                show_date: 'bool',
                date_pub: 'string',
                date_end: 'string',
                image: 'string',
                hits: 'int',
                sort: 'int',
                stars: 'float',
                path: 'string',
                show_image: 'bool',
                name: 'string',

                // теги/ресурсы/инфо-блоки/галереи
                post_tag_id: 'int',
                resource_uuid: 'string',
                position: 'string',
                info_block_id: 'int',
                gallery_id: 'int',
                original_name: 'string',
                file: 'string',

                // аналитика
                request_uuid: 'string',
                timestamp: 'string',
                method: 'string',
                host: 'string',
                query: 'string',
                latency: 'int',
                ip: 'string',
                country: 'string',
                city: 'string',
                browser: 'string',
                device: 'string',
                os: 'string',
                language: 'string',
                referrer: 'string',
                resolution_width: 'int',
                resolution_height: 'int',
                request_size: 'int',
                response_size: 'int',
                utm_campaign: 'string',
                utm_source: 'string',
                utm_medium: 'string',

                // сообщения
                from: 'string',
                to: 'string',
                subject: 'string',
                body: 'string',
                attachment: 'string',
                viewed: 'bool',

                // меню
                ico: 'string',
                publisher_uuid: 'string',
                menu_id: 'int',
                menu_item_id: 'int',

                // файлы
                size: 'int',
                type: 'string',
                is_received: 'bool',
            });
        }

    },
    resolution: function () {
        const cookie = new this.cookie('resolution');
        if (!cookie.get()) {
            cookie.value = window.screen.width + ";" + window.screen.height;
            cookie.options = {expires: 86400, path: '/'};
            cookie.set();
        }
    },
    uuid: function () {
        return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function (c) {
            const r = Math.random() * 16 | 0,
                v = c === 'x' ? r : (r & 0x3 | 0x8);
            return v.toString(16);
        });
    },
    timeUUID: function () {
        return Date.now().toString(36) + Math.random().toString(36).substr(2);
    },
    phone: function (s, plus = true) {
        const startsWith = plus ? '+7' : '8';
        let phone = s.replace(/[^0-9]/g, '');
        if (phone.startsWith('7') && plus) {
            phone = phone.substr(1);
        }
        if (phone.startsWith('8')) {
            phone = phone.substr(1);
        }
        return phone.replace(/(\d{3})(\d{3})(\d{2})(\d{2})/g, `${startsWith} ($1) $2 $3 $4`);
    },
    price: function (_number) {
        const decimal = 0;
        const separator = ' ';
        const decpoint = '.';
        const format_string = '# ₽';
        let r = parseFloat(_number);
        const exp10 = Math.pow(10, decimal);
        r = Math.round(r * exp10) / exp10;
        let rr = Number(r).toFixed(decimal).toString().split('.');
        let b = rr[0].replace(/(\d{1,3}(?=(\d{3})+(?:\.\d|\b)))/g, "\$1" + separator);
        r = (rr[1] ? b + decpoint + rr[1] : b);
        return format_string.replace('#', r);
    },
    inputMask: function (selector) {
        const obj = $(selector);
        if (obj.length) {
            obj.inputmask({"mask": "+7(999) 999-99-99"});
        }
    },
    synchronization: function () {
        const self = this;
        $('body').on('change', '[data-synchronization]', function (evt) {
            let field = $(this);
            let value = field.val();
            let name = field.attr('data-synchronization').split('.');
            name.forEach(function (item, i, arr) {
                let selector = `[name="${item}"]`;
                $(selector).val(value);
            });
        })
    },
    lazyLoading: {
        images: [],
        selector: null,
        loading: function (target, attribute) {
            const self = this;
            this.selector = target;
            const blocks = $(target);
            if (!blocks.length) {
                return;
            }
            this.start(blocks, attribute);
            const _window = $(window);
            _window.scroll(function () {
                const _top = _window.scrollTop();
                const _height = _window.height();
                self.images.forEach(function (item, index, object) {
                    if (_top + _height >= $(item).offset().top) {
                        const atr = target.replace(/[\.\#\[\]]/gi, '');
                        $(item).attr(attribute, $(item).attr(atr));
                        item.removeAttribute(atr);
                        object.splice(index, 1);
                    }
                });
            });
        },
        start: function (blocks, attribute) {
            const _window = $(window);
            const _top = _window.scrollTop();
            const _height = _window.height();
            for (let val of blocks) {
                if (_top + _height >= $(val).offset().top) {
                    const atr = this.selector.replace(/[\.\#\[\]]/gi, '');
                    $(val).attr(attribute, $(val).attr(atr));
                    val.removeAttribute(atr);
                } else {
                    this.images.push(val);
                }
            }
        },
    },
    isEmpty: function (value) {
        return (
            value === null
            || value === undefined
            || value === '0'
            || value === 'false'
            || value === false
            || (typeof value === 'string' && value.trim() === '')
            || (Array.isArray(value) && value.length === 0)
        );
    },
    run: function () {
        try {
            const urlSearchParams = new URLSearchParams(window.location.search);
            const params = Object.fromEntries(urlSearchParams.entries());
            if (Object.keys(params).length) {
                this.pathSearchParams = params;
            }
            const path = document.location.pathname.replace(/\//, '');
            if (path) {
                this.path = path;
                this.pathArray = path.split('/');
            }
            const hash = document.location.hash.replace(/\#/, '');
            if (hash) {
                this.pathHash = hash;
            }
        } catch (e) {
            this.console.error(e.message);
        }
        try {
            this.inputMask('.phone-mask');
        } catch (e) {
            this.console.error(e.message);
        }
        this.validation.control();
        this.setMaps();
        this.resolution();
        this.synchronization();
        this.lazyLoading.loading('[data-js-image-lazy-loading]', 'src');
        this.typeHints.run()
    }
}