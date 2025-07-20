import { createRoot, createSignal } from 'solid-js';
import { Second } from '../constants';

export type NotificationProps = {
    message: string;
    title?: string;
    duration?: number;
};

export type Notification = {
    id: number;
    title: string;
    message?: string;
    type: 'info' | 'success' | 'warning' | 'error' | 'loading' | 'aurevoir';
    duration: number;
};

export const useNotification = createRoot(() => {
    const [notifications, setNotifications] = createSignal<Notification[]>([]);

    const removeNotification = (id: number) => setNotifications(prevNotifications => prevNotifications.filter(notification => notification.id !== id));

    const e = (props: NotificationProps) => {
        const notification: Notification = {
            id: Date.now(),
            title: props.title ?? 'an error occurred',
            message: props.message,
            type: 'error',
            duration: props.duration ?? 5 * Second,
        };
        setNotifications(notifications => [...notifications, notification]);
        setTimeout(() => removeNotification(notification.id), notification.duration);
        return notification.id;
    };

    const i = (props: NotificationProps) => {
        const notification: Notification = {
            id: Date.now(),
            title: props.title ?? 'info',
            message: props.message,
            type: 'info',
            duration: props.duration ?? 5 * Second,
        };
        setNotifications(notifications => [...notifications, notification]);
        setTimeout(() => removeNotification(notification.id), notification.duration);
        return notification.id;
    };

    const s = (props: NotificationProps) => {
        const notification: Notification = {
            id: Date.now(),
            title: props.title ?? 'success',
            message: props.message,
            type: 'success',
            duration: props.duration ?? 3 * Second,
        };
        setNotifications(notifications => [...notifications, notification]);
        setTimeout(() => removeNotification(notification.id), notification.duration);
        return notification.id;
    };

    const w = (props: NotificationProps) => {
        const notification: Notification = {
            id: Date.now(),
            title: props.title ?? 'warning',
            message: props.message,
            type: 'warning',
            duration: props.duration ?? 5 * Second,
        };
        setNotifications(notifications => [...notifications, notification]);
        setTimeout(() => removeNotification(notification.id), notification.duration);
        return notification.id;
    };

    const l = (props: NotificationProps) => {
        const notification: Notification = {
            id: Date.now(),
            title: props.title ?? 'loading',
            message: props.message,
            type: 'loading',
            duration: props.duration ?? 1 * Second,
        };
        setNotifications(notifications => [...notifications, notification]);
        setTimeout(() => removeNotification(notification.id), notification.duration);
        return notification.id;
    };

    const aurevoir = (props: NotificationProps) => {
        const notification: Notification = {
            id: Date.now(),
            title: props.title ?? 'Bye',
            message: props.message,
            type: 'aurevoir',
            duration: props.duration ?? 3 * Second,
        };
        setNotifications(notifications => [...notifications, notification]);
        setTimeout(() => removeNotification(notification.id), notification.duration);
        return notification.id;
    };

    return { notifications, l, w, aurevoir, s, i, e, removeNotification };
});
