import React from 'react';
import {
    Accordion,
    AccordionDetails,
    AccordionSummary,
    Table,
    TableBody,
    TableCell,
    TableRow,
    Typography
} from '@material-ui/core';
import {
    getPluginConfig
} from '../../api/PluginConfig';
import {
    Cat as CatSubject,
    Goose as GooseSubject,
    PluginConfig as Subject,
    Secret as SecretSubject,
    Size as SizeSubject,
    Welcome as WelcomeSubject
} from '../../models';
import {
    ExpandMore as ExpandMoreIcon
} from '@material-ui/icons';

interface secretProps {
    subject: SecretSubject
}

function Secret({ subject }: secretProps) {
    if (subject == null) {
        return null;
    }
    if (subject.valueFrom) {
        return (
            <>
                <Typography>Secret</Typography>
                <Typography>{subject.valueFrom.secretKeyRef.name}</Typography>
                <Typography>Key</Typography>
                <Typography>{subject.valueFrom.secretKeyRef.key}</Typography>
            </>
        )
    }
    return (
        <>
            <Typography>Key</Typography>
            <Typography>{subject.value}</Typography>
        </>
    )
}

interface catProps {
    subject: CatSubject
}

function Cat({ subject }: catProps) {
    if (subject == null) {
        return null;
    }
    return (
        <Accordion>
            <AccordionSummary expandIcon={<ExpandMoreIcon />}>
                <Typography>Cat</Typography>
            </AccordionSummary>
            <AccordionDetails>
                <Secret subject={subject.key} />
            </AccordionDetails>
        </Accordion>
    )
}

interface gooseProps {
    subject: GooseSubject
}

function Goose({ subject }: gooseProps) {
    if (subject == null) {
        return null;
    }
    return (
        <Accordion>
            <AccordionSummary expandIcon={<ExpandMoreIcon />}>
                <Typography>Goose</Typography>
            </AccordionSummary>
            <AccordionDetails>
                <Secret subject={subject.key} />
            </AccordionDetails>
        </Accordion>
    )
}

interface sizeProps {
    subject: SizeSubject
}

function Size({ subject }: sizeProps) {
    if (subject == null) {
        return null;
    }
    return (
        <Accordion>
            <AccordionSummary expandIcon={<ExpandMoreIcon />}>
                <Typography>Size</Typography>
            </AccordionSummary>
            <AccordionDetails>
                <Table>
                    <TableBody>
                        <TableRow>
                            <TableCell>Small</TableCell>
                            <TableCell>{subject.s}</TableCell>
                        </TableRow>
                        <TableRow>
                            <TableCell>Medium</TableCell>
                            <TableCell>{subject.m}</TableCell>
                        </TableRow>
                        <TableRow>
                            <TableCell>Large</TableCell>
                            <TableCell>{subject.l}</TableCell>
                        </TableRow>
                        <TableRow>
                            <TableCell>Extra large</TableCell>
                            <TableCell>{subject.xl}</TableCell>
                        </TableRow>
                        <TableRow>
                            <TableCell>Extra extra large</TableCell>
                            <TableCell>{subject.xxl}</TableCell>
                        </TableRow>
                    </TableBody>
                </Table>
            </AccordionDetails>
        </Accordion >
    )
}

interface welcomeProps {
    subject: WelcomeSubject
}

function Welcome({ subject }: welcomeProps) {
    if (subject == null) {
        return null;
    }
    return (
        <Accordion>
            <AccordionSummary expandIcon={<ExpandMoreIcon />}>
                <Typography>Welcome</Typography>
            </AccordionSummary>
            <AccordionDetails>
                <Typography>Template</Typography>
                <Typography>{subject.messageTemplate}</Typography>
            </AccordionDetails>
        </Accordion>
    )
}

interface props {
    match: any
}

export function PluginConfig({ match }: props) {
    const [subject, setSubject] = React.useState<Subject | null>(null);

    React.useEffect(() => {
        getPluginConfig(match.params.name).then(result => setSubject(result));
    }, [match.params.name]);

    if (subject == null) {
        return null;
    }

    return (
        <>
            <Cat subject={subject.spec.cat}></Cat>
            <Goose subject={subject.spec.goose}></Goose>
            <Size subject={subject.spec.size}></Size>
            <Welcome subject={subject.spec.welcome}></Welcome>
        </>
    );
}
