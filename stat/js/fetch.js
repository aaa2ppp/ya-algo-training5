// ==UserScript==
// @name         New Userscript
// @namespace    panfilov.pp.ru
// @version      2024-03-20
// @description  try to take over the world!
// @author       You
// @match        https://contest.yandex.ru/contest/*/enter/*
// @icon         https://www.google.com/s2/favicons?sz=64&domain=google.com
// @grant        GM_setClipboard
// ==/UserScript==

(function() {
    'use strict';

    if (!document.location.pathname.endsWith('/enter/')) {
        throw new Error(`path must ends with "/enter/", got "${document.location.pathname}`);
    }

    // XXX возможно это избыточно
    if (document.readyState == 'loading') {
        // ещё загружается, ждём события
        document.addEventListener('DOMContentLoaded', main);
    } else {
        // DOM готов!
        main();
    }

    const maxFetchCount = 200;
    const fetchPeriod = 1000;
    const fetchedPages = new Array();

    let stop = false;

    function main() {
        let pathname = document.location.pathname;
        pathname = pathname.substring(0, pathname.length-'/enter/'.length) + '/standings/';

        const url = document.location.protocol + '//' + document.location.hostname + pathname;
        console.log(`url: ${url}`);

        const out = document.createElement('span');
        out.textContent = 'press Start to fetch standings';

        const toolbar = document.createElement('div');
        toolbar.insertAdjacentElement('beforeend', newButton(
            'Start',
            function() {
                fetchedPages.length = 0;
                stop = false;
                fetchStandings(url, 1, out);
            }
        ));

        toolbar.insertAdjacentElement('beforeend', newButton(
            'Stop',
            function() {
                stop = true;
            }
        ));

        toolbar.insertAdjacentElement('beforeend', newButton(
            'Get',
            function() {
                getResult(out);
            }
        ));

        toolbar.insertAdjacentHTML('beforeend', '&nbsp;');
        toolbar.insertAdjacentElement('beforeend', out);

        document.body.insertAdjacentElement('afterbegin', toolbar);
    };

    function newButton(label, onclick) {
        const button = document.createElement('button');
        button.textContent = label;
        button.onclick = onclick;
        return button;
    }

    function getResult(out) {
        GM_setClipboard(fetchedPages.join('\n<!-- EOP -->\n'));
    }

    function fetchStandings(url, pageNum, out) {
        if (stop) {
            out.textContent = 'fetch stopped';
            return;
        }

        const curUrl = `${url}?p=${pageNum}`;
        console.log(`fetch: ${curUrl}`);
        out.textContent = `fetch: ${curUrl}`;

        fetch(curUrl)
            .then((response) => {

            if (!response.ok) {
                out.textContent += ` HTTP error! Status: ${response.status}`
                throw new Error(`HTTP error! Status: ${response.status}`);
            }

            return response.text();
        })
            .then((response) => {

            fetchedPages[fetchedPages.length] = extractStandingsTab(response);

            if (pageNum < maxFetchCount && pageNum < findLastPageNum(response)) {
                setTimeout(function() {fetchStandings(url, pageNum+1, out)}, fetchPeriod);
            } else {
                out.textContent = "fetch done. Press Get to copy data to clipboard";
            }
        });
    }

    function extractStandingsTab(html) {
        const p1 = html.indexOf('<table class="table table_role_standings');
        const p2 = html.indexOf('</table>') + '</table>'.length;
        return html.substring(p1, p2);
    }

    function findLastPageNum(html) {
        let lastPage;
        const s1 = '<span class="button__text">';
        const s2 = '</span>';

        for (let i = 0, p1 = html.length-1; i < 2 && p1 >= 0; i++, p1--) {
            p1 = html.lastIndexOf(s1, p1);
            if (p1 == -1) {
                console.log(`findLastPageNum: not found "${s1}"`);
                return NaN;
            }

            const p2 = html.indexOf(s2, p1);
            if (p2 == -1) {
                console.log(`findLastPageNum: not found "${s2}">`);
                return NaN;
            }

            const nums = html.substring(p1+s1.length, p2);
            console.log(`findLastPageNum: ${nums}`);

            const num = Number(nums);
            if (!isNaN(num)) {
                return num+i;
            }
        }

        return NaN;
    }
})();
