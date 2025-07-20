import { type Accessor, createEffect, createSignal, type Setter } from 'solid-js';
import { createStore, type SetStoreFunction, type Store } from 'solid-js/store';

const localStoragePrefix = 'MEH_';

export function createLocalStore<T extends object>(name: string, init: T): [Store<T>, SetStoreFunction<T>] {
    const localState = localStorage.getItem(localStoragePrefix + name);
    const [state, setState] = createStore<T>(localState ? JSON.parse(localState) : init);
    createEffect(() => localStorage.setItem(localStoragePrefix + name, JSON.stringify(state)));
    return [state, setState];
}

export function createLocalSignal<T extends string | number | boolean | Record<string, unknown> | null>(name: string, init: T): [Accessor<T>, Setter<T>] {
    name = localStoragePrefix + name;
    const localState = localStorage.getItem(name);
    let localStateParsed: T | undefined;

    if (localState) {
        try {
            localStateParsed = JSON.parse(localState) as T;
        } catch (_error) {
            console.error('Error parsing local storage value:', name, localState);
        }
    }

    const [state, setState] = createSignal<T>(localStateParsed ?? init);

    createEffect(() => {
        if (state() === null) {
            localStorage.removeItem(name);
        } else {
            localStorage.setItem(name, JSON.stringify(state()));
        }
    });

    return [state, setState];
}
