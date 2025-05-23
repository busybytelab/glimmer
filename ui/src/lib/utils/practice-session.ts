import type { BreadcrumbItem } from '$lib/types';
import type { SessionWithExpandedData } from '$lib/services/session';

/**
 * Updates breadcrumb items for a practice session
 */
export function updateBreadcrumbs(session: SessionWithExpandedData | null): BreadcrumbItem[] {
    if (!session) return [];
    
    const items: BreadcrumbItem[] = [
        {
            label: 'Topics',
            href: '/practice-topics',
            icon: 'topic'
        }
    ];
    
    if (session.expand?.practice_topic) {
        items.push({
            label: session.expand.practice_topic.name,
            href: `/practice-topics/${session.expand.practice_topic.id}`,
            icon: 'topic'
        });
    }
    
    items.push({
        label: session.name || 'Practice Session',
        icon: 'session'
    });
    
    return items;
}

/**
 * Handles print functionality for practice sessions
 */
export function handlePrint(): void {
    let printMode = true;
    setTimeout(() => {
        window.print();
        setTimeout(() => {
            printMode = false;
        }, 500);
    }, 200);
} 