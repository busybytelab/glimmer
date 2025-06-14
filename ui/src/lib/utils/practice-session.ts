import type { BreadcrumbItem, BreadcrumbIcon } from '$lib/types';
import type { SessionWithExpandedData } from '$lib/services/session';

/**
 * Updates breadcrumb items for a practice session
 */
export function updateBreadcrumbs(session: SessionWithExpandedData | null): BreadcrumbItem[] {
    if (!session) return [];
    
    // Check if we're in a learner route by looking at the current URL
    const isLearnerRoute = window.location.pathname.startsWith('/learners/');
    const learnerId = isLearnerRoute ? window.location.pathname.split('/')[2] : null;
    
    const items: BreadcrumbItem[] = [
        {
            label: 'Topics',
            href: isLearnerRoute 
                ? `/learners/${learnerId}/practice-topics`
                : '/account/practice-topics',
            icon: 'topic' as BreadcrumbIcon
        }
    ];
    
    if (session.expand?.practice_topic) {
        items.push({
            label: session.expand.practice_topic.name,
            href: isLearnerRoute
                ? `/learners/${learnerId}/practice-topics/${session.expand.practice_topic.id}`
                : `/account/practice-topics/${session.expand.practice_topic.id}`,
            icon: 'topic' as BreadcrumbIcon
        });
    }
    
    items.push({
        label: session.name || 'Practice Session',
        icon: 'session' as BreadcrumbIcon
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