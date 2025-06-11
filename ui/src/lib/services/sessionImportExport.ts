import pb from '$lib/pocketbase';
import type { ExportedSession, PracticeSession, PracticeTopic } from '$lib/types';
import { sessionService } from './session';
import { accountService } from './accounts';

/**
 * Service for handling practice session import/export functionality
 * Follows the single responsibility principle by focusing only on import/export operations
 */
class SessionImportExportService {
    /**
     * Exports a practice session with its associated data (topic, items)
     * Excludes user-specific data to allow session reuse
     */
    async exportPracticeSession(sessionId: string): Promise<ExportedSession> {
        try {
            // Load the complete session with expanded data
            const session = await sessionService.loadSession(sessionId);
            
            if (!session.expand?.practice_topic || !session.expand?.practice_items) {
                throw new Error('Failed to load session data: Missing topic or practice items');
            }

            const topic = session.expand.practice_topic as PracticeTopic;

            // Create the export data structure
            const exportData: ExportedSession = {
                session: {
                    id: session.id,
                    created: session.created,
                    updated: session.updated,
                    collectionId: session.collectionId,
                    collectionName: session.collectionName,
                    name: session.name,
                    status: session.status,
                    assigned_at: session.assigned_at,
                    completed_at: session.completed_at,
                    generation_prompt: session.generation_prompt,
                    practice_topic: session.practice_topic,
                    practice_items: session.practice_items,
                    score: session.score
                },
                topic: {
                    id: topic.id,
                    name: topic.name,
                    subject: topic.subject,
                    description: topic.description,
                    target_age_range: topic.target_age_range,
                    target_grade_level: topic.target_grade_level,
                    learning_goals: topic.learning_goals,
                    base_prompt: topic.base_prompt,
                    system_prompt: topic.system_prompt,
                    tags: topic.tags,
                    llm_model: topic.llm_model
                },
                practice_items: session.expand.practice_items.map(item => ({
                    id: item.id,
                    created: item.created,
                    updated: item.updated,
                    collectionId: item.collectionId,
                    collectionName: item.collectionName,
                    question_text: item.question_text,
                    question_type: item.question_type,
                    options: item.options,
                    correct_answer: item.correct_answer,
                    explanation: item.explanation,
                    explanation_for_incorrect: item.explanation_for_incorrect,
                    hints: item.hints,
                    difficulty_level: item.difficulty_level,
                    status: item.status,
                    tags: item.tags,
                    practice_topic: item.practice_topic,
                    review_status: item.review_status,
                    review_date: item.review_date
                }))
            };

            return exportData;
        } catch (error) {
            console.error('Failed to export session:', error);
            throw new Error('Failed to export practice session');
        }
    }

    /**
     * Imports a practice session for a specific learner
     * Creates or reuses the topic and creates new practice items
     */
    async importPracticeSession(
        importData: ExportedSession, 
        learnerId: string,
        existingTopicId: string | null = null
    ): Promise<PracticeSession> {
        try {
            // Get the current user's account
            const account = await accountService.getAccount();
            if (!account?.id) {
                throw new Error('No account found for current user');
            }

            console.log('Starting import process with account:', account.id);

            // If existingTopicId is provided, use that instead of creating a new topic
            const topicId = existingTopicId || (await pb.collection('practice_topics').create({
                ...importData.topic,
                id: undefined,
                created: undefined,
                updated: undefined,
                collectionId: undefined,
                collectionName: undefined,
                account: account.id
            })).id;

            console.log('Using topic ID:', topicId);

            // Create practice items sequentially to avoid request cancellation
            const practiceItems = [];
            for (const item of importData.practice_items) {
                try {
                    const createdItem = await pb.collection('practice_items').create({
                        question_text: item.question_text,
                        question_type: item.question_type || 'multiple_choice',
                        options: JSON.stringify(item.options || {}),
                        correct_answer: item.correct_answer,
                        explanation: item.explanation,
                        explanation_for_incorrect: JSON.stringify(item.explanation_for_incorrect || {}),
                        hints: JSON.stringify(item.hints || []),
                        difficulty_level: item.difficulty_level || 'medium',
                        status: 'Imported',
                        tags: JSON.stringify(item.tags || {}),
                        practice_topic: topicId,
                        account: account.id,
                        review_status: item.review_status,
                        review_date: item.review_date
                    });
                    practiceItems.push(createdItem);
                    console.log('Created practice item:', createdItem.id);
                } catch (error) {
                    console.error('Failed to create practice item:', error);
                    throw error;
                }
            }

            console.log('All practice items created:', practiceItems.map(item => item.id));

            // Create new session with the selected topic ID and practice items
            const sessionData = {
                name: importData.session.name || 'Imported Session',
                status: 'Imported',
                assigned_at: new Date().toISOString(),
                generation_prompt: importData.session.generation_prompt || '',
                learner: learnerId,
                practice_topic: topicId,
                account: account.id,
                practice_items: JSON.stringify(practiceItems.map(item => item.id))
            };

            console.log('Creating session with data:', sessionData);
            
            try {
                const session = await pb.collection('practice_sessions').create(sessionData);
                console.log('Session created successfully:', session.id);
                return session;
            } catch (error) {
                console.error('Failed to create session:', error);
                // Clean up created practice items if session creation fails
                for (const item of practiceItems) {
                    try {
                        await pb.collection('practice_items').delete(item.id);
                    } catch (cleanupError) {
                        console.error('Failed to clean up practice item:', cleanupError);
                    }
                }
                throw error;
            }
        } catch (error) {
            console.error('Failed to import session:', error);
            throw error;
        }
    }

    /**
     * Downloads the exported session as a JSON file
     */
    downloadExportedSession(exportData: ExportedSession, filename?: string): void {
        const blob = new Blob([JSON.stringify(exportData, null, 2)], { type: 'application/json' });
        const url = URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;

        // Create a filename that includes both topic name and session ID
        const topicName = exportData.topic.name.toLowerCase().replace(/[^a-z0-9]+/g, '_');
        const defaultFilename = `practice_session_${topicName}_${exportData.session.id}.json`;
        a.download = filename || defaultFilename;

        document.body.appendChild(a);
        a.click();
        document.body.removeChild(a);
        URL.revokeObjectURL(url);
    }
}

export const sessionImportExportService = new SessionImportExportService(); 