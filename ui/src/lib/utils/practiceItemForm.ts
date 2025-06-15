import type { PracticeItem } from '$lib/types';

/**
 * Interface for form data structure
 */
export interface PracticeItemFormData {
    questionText: string;
    correctAnswer: string;
    explanation: string;
    hints: string;
    options: string;
    incorrectAnswers: Array<{ answer: string; explanation: string }>;
    tags: string;
}

/**
 * Convert PracticeItem to form data structure
 */
export function practiceItemToFormData(item: PracticeItem): PracticeItemFormData {
    // Handle incorrect answers
    let incorrectAnswers: Array<{ answer: string; explanation: string }> = [];
    if (item.explanation_for_incorrect) {
        try {
            incorrectAnswers = Object.entries(item.explanation_for_incorrect).map(([answer, explanation]) => ({
                answer,
                explanation
            }));
        } catch (err) {
            console.error('Failed to parse explanation_for_incorrect:', err);
            incorrectAnswers = [];
        }
    }

    return {
        questionText: item.question_text,
        correctAnswer: typeof item.correct_answer === 'string' 
            ? item.correct_answer 
            : JSON.stringify(item.correct_answer),
        explanation: item.explanation,
        hints: Array.isArray(item.hints) ? item.hints.join('\n') : '',
        options: Array.isArray(item.options) ? item.options.join('\n') : '',
        incorrectAnswers,
        tags: Array.isArray(item.tags) ? item.tags.join(', ') : ''
    };
}

/**
 * Convert form data to PracticeItem update structure
 */
export function formDataToPracticeItemUpdate(formData: PracticeItemFormData): Partial<PracticeItem> {
    // Parse hints back into array
    const hintsArray = formData.hints
        .split('\n')
        .map(hint => hint.trim())
        .filter(hint => hint);
    
    // Parse options back into array
    const optionsArray = formData.options
        .split('\n')
        .map(option => option.trim())
        .filter(option => option);
    
    // Convert incorrect answers to object
    const explanationForIncorrect = formData.incorrectAnswers.reduce((acc, { answer, explanation }) => {
        if (answer.trim() && explanation.trim()) {
            acc[answer.trim()] = explanation.trim();
        }
        return acc;
    }, {} as Record<string, string>);
    
    // Parse tags back into array
    const tagsArray = formData.tags
        .split(',')
        .map(tag => tag.trim())
        .filter(tag => tag);

    // Return the update object
    return {
        question_text: formData.questionText,
        correct_answer: formData.correctAnswer,
        explanation: formData.explanation,
        hints: hintsArray,
        options: optionsArray,
        explanation_for_incorrect: explanationForIncorrect,
        tags: tagsArray,
        review_status: 'APPROVED',
        review_date: new Date().toISOString()
    };
}

/**
 * Validate form data
 */
export function validatePracticeItemForm(formData: PracticeItemFormData): string[] {
    const errors: string[] = [];

    if (!formData.questionText.trim()) {
        errors.push('Question text is required');
    }

    if (!formData.correctAnswer.trim()) {
        errors.push('Correct answer is required');
    }

    if (!formData.explanation.trim()) {
        errors.push('Explanation is required');
    }

    // Validate incorrect answers
    formData.incorrectAnswers.forEach((incorrect, index) => {
        if (incorrect.answer.trim() && !incorrect.explanation.trim()) {
            errors.push(`Explanation for incorrect answer #${index + 1} is required`);
        }
        if (!incorrect.answer.trim() && incorrect.explanation.trim()) {
            errors.push(`Incorrect answer #${index + 1} text is required`);
        }
    });

    return errors;
} 